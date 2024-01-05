package main

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sshModel struct {
	client  *Client
	input   textinput.Model
	output  string
	spinner spinner.Model
	err     error
}

type errMsg error

// pass host here (username and host are currently hardcoded in .env)
func initSSHModel() sshModel {
	// ssh
	signer, err := LoadPrivKey()
	if err != nil {
		log.Fatal("failed to laod private key: ", err)
	}
	client := NewSSHClient(signer)

	// input
	t := textinput.New()
	t.Placeholder = "Enter command"
	t.Focus()

	// spinner
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return sshModel{
		client:  client,
		input:   t,
		spinner: s,
	}
}

func (m sshModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m sshModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "X":
			output, err := m.client.RunCmd("ls -l")
			if err != nil {
				log.Fatal("failed to run command: ", err)
			}
			m.output = output
			return m, nil

		case "q", "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			input := m.input.Value()
			output, err := m.client.RunCmd(input)

			m.output = output
			m.input.SetValue("")

			if err != nil {
				m.err = errMsg(err)
				return m, nil
			}

			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	// input
	m.input, cmd = m.input.Update(msg)

	// spinner
	var cmdSpinner tea.Cmd
	m.spinner, cmdSpinner = m.spinner.Update(msg)
	return m, tea.Batch(cmd, cmdSpinner)
}

func (m sshModel) View() string {

	var b strings.Builder
	b.WriteString("SSEX\n")
	b.WriteString(m.input.View())
	b.WriteString("\n")
	b.WriteString(m.output)
	b.WriteString("\n")
	b.WriteString(m.spinner.View())
	return b.String()
}
