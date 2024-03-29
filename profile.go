package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type state int

type sshResult struct {
	client *Client
	err    error
}

const (
	formState state = iota
	connectingState
)

const (
	DialTimeout = 15
)

type connectionEstablishedMsg struct {
	client *Client
}

type Profile struct {
	Host         string
	Port         string
	User         string
	IdentityFile string
}

type profileModel struct {
	client  *Client
	form    *huh.Form // huh.Form is a tea.Model
	profile Profile
	state   state
	spinner spinner.Model
	height  int
	width   int
	help    help.Model
	keys    KeyMap
	err     error
}

func initForm(p *Profile) *huh.Form {
	// form
	f := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Value(&p.Host).
				Title("Host").
				Key("Host").
				Placeholder("e.g 192.149.252.76").
				Validate(func(s string) error {
					if s == "" {
						return errors.New("host is required")
					}
					return nil
				}).
				Description("Enter the hostname or IP address of the remote server:").
				Prompt("$ "),
			huh.NewInput().
				Value(&p.Port).
				Title("Port").
				Key("Port").
				Placeholder("e.g 12345").
				Validate(func(s string) error {
					if s != "" && !IsNumber(s) {
						return errors.New("port must be a number")
					}
					return nil
				}).
				Description("Port to connect to on the remote host:").
				Prompt("$ "),
			huh.NewInput().
				Value(&p.User).
				Title("User").
				Key("User").
				Placeholder("e.g root").
				Validate(func(s string) error {
					if s == "" {
						return errors.New("user is required")
					}
					return nil
				}).
				Description("Specifies the user to log in as on the remote machine:").
				Prompt("$ "),
		),
	)
	f.WithTheme(huh.ThemeCatppuccin())
	f.WithShowHelp(false)
	return f
}

func initSpinner() spinner.Model {
	// spinner
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(c100))
	return s
}
func initialModel() profileModel {
	var p Profile

	f := initForm(&p)
	s := initSpinner()
	k := DefaultKeyMap
	k.New.Unbind()

	return profileModel{
		form:    f,
		profile: p,
		state:   formState,
		spinner: s,
		help:    help.New(),
		keys:    DefaultKeyMap,
	}

}

func (m profileModel) Init() tea.Cmd {
	return tea.Batch(
		m.form.Init(),
		tea.EnterAltScreen,
	)
}

func (m profileModel) reset() (tea.Model, tea.Cmd) {
	profile := Profile{}
	m.form = initForm(&profile)
	m.profile = profile
	m.spinner = initSpinner()
	m.state = formState
	m.keys = DefaultKeyMap
	m.err = nil
	return m, m.form.Init()
}

func (m profileModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.New):
			return m.reset()
		}
	case connectionEstablishedMsg:
		var newModel sshModel
		return initSSHModel(msg.client, m.profile, m.height, m.width), newModel.Init()
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case errMsg:
		m.err = msg
		m.keys.New.SetEnabled(true)
		return m, nil
	default:
		// completed
		if m.form.State == huh.StateCompleted {
			m.profile = Profile{
				Host: m.form.GetString("Host"),
				Port: m.form.GetString("Port"),
				User: m.form.GetString("User"),
			}

			m.state = connectingState
			m.keys.Back.Unbind()
			m.keys.Next.Unbind()
			m.keys.Clear.Unbind()
			return m, tea.Batch(sshCmd(m.profile), m.spinner.Tick)
		}
	}

	return m, cmd
}

func (m profileModel) View() string {
	var b strings.Builder

	// title
	title := buildTitle(m.width)
	b.WriteString(title)
	b.WriteString("\n")

	// submited, connecting
	if m.form.State == huh.StateCompleted {
		if m.err != nil {
			errText := m.err.Error()

			// failed connection
			styleConnectionError.Width(m.width).ColorWhitespace(false)
			paddingLen := horizontalPadLength(errText, m.width)
			styleConnectionError.PaddingLeft(paddingLen)
			b.WriteString(styleConnectionError.Render(errText))
		} else {
			// connecting...
			connectingStr := buildConnectingStr(m.form, m.spinner, m.width)
			b.WriteString(connectingStr)
			b.WriteString("\n")
		}
	}

	// form
	b.WriteString(m.form.View())
	// empty lines
	buildEmptyLine(&b, m.height)
	// help
	b.WriteString(m.help.View(m.keys))
	return styleApp.Render(b.String())
}

func sshCmd(profile Profile) tea.Cmd {
	return func() tea.Msg {
		// load private key
		signer, err := LoadPrivKey()
		if err != nil {
			return errMsg(
				fmt.Errorf("failed to load private key: %v", err),
			)
		}

		// ssh timeout
		ctx, cancel := context.WithTimeout(context.Background(), DialTimeout*time.Second)
		defer cancel()

		// create ssh client
		resultChan := make(chan sshResult)
		go NewSSHClient(ctx, signer, profile, resultChan)

		// wait for result
		select {
		case result := <-resultChan:
			if result.err != nil {
				return errMsg(
					fmt.Errorf("failed to create ssh client: %v", result.err),
				)
			} else {
				return connectionEstablishedMsg{
					client: result.client,
				}
			}
		case <-ctx.Done():
			return errMsg(fmt.Errorf("🔥 connection timeout"))
		}
	}
}

func buildEmptyLine(b *strings.Builder, height int) {
	contentHeight := lipgloss.Height(b.String())
	if contentHeight < height {
		b.WriteString(strings.Repeat("\n", height-contentHeight))
	}
}

func buildTitle(width int) string {
	b := strings.Builder{}
	styleChar := lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Bold(true)

	b.WriteString("\n")
	if width > 0 { // width is 0 when the app is starting
		padding := strings.Repeat(" ", horizontalPadLength("SSEX", width))
		b.WriteString(padding)
	}
	b.WriteString(styleChar.Background(lipgloss.Color(c500)).Render("S"))
	b.WriteString(styleChar.Background(lipgloss.Color(c600)).Render("S"))
	b.WriteString(styleChar.Background(lipgloss.Color(c700)).Render("E"))
	b.WriteString(styleChar.Background(lipgloss.Color(c800)).Render("X"))
	b.WriteString("\n")

	return b.String()
}

func buildConnectingStr(f *huh.Form, s spinner.Model, width int) string {
	host := f.GetString("Host")
	port := f.GetString("Port")
	user := f.GetString("User")

	if port == "" {
		// only for display purpose
		port = "22"
	}
	hostRender := styleHost.Render(host)
	portRender := stylePort.Render(port)
	userRender := styleUser.Render(user)

	// connecting
	connectingStr := fmt.Sprintf(
		"%s Connecting to %s:%s as %s...",
		s.View(),
		hostRender,
		portRender,
		userRender,
	)

	var styleConnectingStr lipgloss.Style
	// pad left
	if width > 0 {
		paddingLen := horizontalPadLength(connectingStr, width)
		styleConnectingStr = lipgloss.NewStyle().PaddingLeft(paddingLen)
	}
	return styleConnectingStr.Render(connectingStr)
}

func horizontalPadLength(s string, width int) int {
	l := (width / 2) - (lipgloss.Width(s) / 2) - 1 // -1 is not really needed, but cleaner
	return l
}
