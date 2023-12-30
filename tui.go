package main

import (
	"log"

	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	client *Client
	input  textinput.Model
	output string
	err    error
}

type errMsg error

func initialModel() model {
	signer, err := LoadPrivKey()
	if err != nil {
		log.Fatal("failed to laod private key: ", err)
	}
	client := NewSSHClient(signer)
	t := textinput.New()
	t.Placeholder = "Enter command"
	t.Focus()

	return model{
		client: client,
		input:  t,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	return m, cmd
}

func (m model) View() string {

	var b strings.Builder
	b.WriteString("SSEX\n")
	b.WriteString(m.input.View())
	b.WriteString("\n")
	b.WriteString(m.output)
	return b.String()
}
