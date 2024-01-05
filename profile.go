package main

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
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
				Placeholder("e.g github.com").
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
	if m.form.State == huh.StateCompleted {
		return fmt.Sprintf(
			"%s Connecting to %s:%s as %s...",
			m.spinner.View(),
			m.form.GetString("Host"),
			m.form.GetString("Port"),
			m.form.GetString("User"),
		)
	}
	return m.form.View()
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
		client := NewSSHClient(signer, profile.User, address)
		return connectionEstablishedMsg{
			client: client,
		}
	}
}
