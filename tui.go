package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	client *Client
	output string
}

func initialModel() model {
	signer, err := LoadPrivKey()
	if err != nil {
		log.Fatal("failed to laod private key: ", err)
	}
	client := NewSSHClient(signer)
	return model{
		client: client,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "x":
			output, err := m.client.RunCmd("ls -l")
			if err != nil {
				log.Fatal("failed to run command: ", err)
			}
			m.output = output
			return m, nil

		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	var s string
	if m.output != "" {
		s += m.output
	} else {
		s += "its pretty lonely out here.."
	}
	return s
}
