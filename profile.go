package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type connectionEstablishedMsg struct {
	client *Client
}

type connectionErrorMsg error

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
	spinner spinner.Model
	height  int
	width  int
	help    help.Model
	keys    KeyMap
	err     error
}

func initialModel() profileModel {
	var p Profile

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

	// spinner
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(c100))

	return profileModel{
		form:    f,
		profile: p,
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
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case connectionEstablishedMsg:
		var newModel sshModel
		return initSSHModel(msg.client, m.height, m.width), newModel.Init()
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case connectionErrorMsg:
		m.err = msg
		// ask user to press something to restart form
	}

	// completed
	if m.form.State == huh.StateCompleted {
		// uncessary, since we only use m.profile here
		m.profile = Profile{
			Host: m.form.GetString("Host"),
			Port: m.form.GetString("Port"),
			User: m.form.GetString("User"),
		}
		return m, tea.Batch(sshCmd(m.profile), m.spinner.Tick)
	}

	// spinner
	var cmdSpinner tea.Cmd
	m.spinner, cmdSpinner = m.spinner.Update(msg)
	return m, tea.Batch(cmd, cmdSpinner)
}

func (m profileModel) View() string {
	var b strings.Builder

	title := buildTitle(m.width)
	b.WriteString(title)

	// submited, connecting
	if m.form.State == huh.StateCompleted {
		if m.err != nil {
			b.WriteString(styleConnectionError.Render(m.err.Error()))
		}

		host := m.form.GetString("Host")
		port := m.form.GetString("Port")
		user := m.form.GetString("User")

		if port == "" {
			// only for display purpose
			port = "22"
		}
		hostRender := styleHost.Render(host)
		portRender := stylePort.Render(port)
		userRender := styleUser.Render(user)

		b.WriteString(
			fmt.Sprintf(
				"%s Connecting to %s:%s as %s...",
				m.spinner.View(),
				hostRender,
				portRender,
				userRender,
			),
		)
		b.WriteString("\n")

		// empty lines
		buildEmptyLine(&b, m.height)
		// help
		b.WriteString(m.help.View(m.keys))
		return styleApp.Render(b.String())
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
		signer, err := LoadPrivKey()
		if err != nil {
			return errMsg(
				fmt.Errorf("failed to load private key: %v", err),
			)
		}

		client, err := NewSSHClient(signer, profile.User, profile.Host, profile.Port)
		if err != nil {
			return errMsg(
				fmt.Errorf("failed to create ssh client: %v", err),
			)
		}

		return connectionEstablishedMsg{
			client: client,
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
    if width > 0 {
        // width is 0 when the app is starting
        padding := strings.Repeat(" ", (width/2) - 5)
        b.WriteString(padding)
    }
	b.WriteString(styleChar.Background(lipgloss.Color(c500)).Render("S"))
	b.WriteString(styleChar.Background(lipgloss.Color(c600)).Render("S"))
	b.WriteString(styleChar.Background(lipgloss.Color(c700)).Render("E"))
	b.WriteString(styleChar.Background(lipgloss.Color(c800)).Render("X"))
	b.WriteString("\n\n")

	return b.String()
}
