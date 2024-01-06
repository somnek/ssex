package main

import (
	"errors"
	"fmt"
	"strings"

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
	form    *huh.Form // huh.Form is a tea.Model
	err     error
	profile Profile
	spinner spinner.Model
	client  *Client
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
				Description("Enter the hostname or IP address of the remote server:"),
			huh.NewInput().
				Value(&p.Port).
				Title("Port").
				Key("Port").
				Placeholder("e.g 12345").
				Validate(func(s string) error {
					if s == "" {
						return errors.New("port is required")
					} else if !IsNumber(s) {
						return errors.New("port must be a number")
					}
					return nil
				}).
				Description("Port to connect to on the remote host:"),
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
				Description("Specifies the user to log in as on the remote machine:"),
		),
	)

	// f.WithTheme()

	// spinner
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))

	return profileModel{
		form:    f,
		profile: p,
		spinner: s,
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
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case connectionEstablishedMsg:
		var newModel sshModel
		return initSSHModel(msg.client), newModel.Init()
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
	var sb strings.Builder

	title := buildTitle()
	sb.WriteString(title)

	if m.form.State == huh.StateCompleted {
		if m.err != nil {
			sb.WriteString(styleConnectionError.Render(m.err.Error()))
		}

		hostRender := styleHost.Render(m.form.GetString("Host"))
		portRender := stylePort.Render(m.form.GetString("Port"))
		userRender := styleUser.Render(m.form.GetString("User"))

		sb.WriteString(
			fmt.Sprintf(
				"%s Connecting to %s:%s as %s...",
				m.spinner.View(),
				hostRender,
				portRender,
				userRender,
			),
		)
		return styleApp.Render(sb.String())
	}

	return styleApp.Render(sb.String() + m.form.View())
}

func sshCmd(profile Profile) tea.Cmd {
	return func() tea.Msg {
		signer, err := LoadPrivKey()
		if err != nil {
			return errMsg(
				fmt.Errorf("failed to load private key: %v", err),
			)
		}

		address := fmt.Sprintf("%s:%s", profile.Host, profile.Port)
		client, err := NewSSHClient(signer, profile.User, address)
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

func buildTitle() string {
	sb := strings.Builder{}
	styleChar := lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Bold(true)
	styleLine := lipgloss.NewStyle().Foreground(lipgloss.Color("6"))

	sb.WriteString(strings.Repeat(" ", 25))
	sb.WriteString(styleChar.Background(lipgloss.Color(c500)).Render("S"))
	sb.WriteString(styleChar.Background(lipgloss.Color(c600)).Render("S"))
	sb.WriteString(styleChar.Background(lipgloss.Color(c700)).Render("H"))
	sb.WriteString(styleChar.Background(lipgloss.Color(c800)).Render("E"))
	sb.WriteString(styleChar.Background(lipgloss.Color(c900)).Render("X"))
	sb.WriteString("\n")
	sb.WriteString(styleLine.Render(strings.Repeat("_", 56)))
	sb.WriteString("\n\n")

	return sb.String()
}
