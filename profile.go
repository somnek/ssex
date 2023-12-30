package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevinburke/ssh_config"
)

type profileModel struct {
	hosts  []*ssh_config.Host
	cursor int
	err    error
}

func initialModel() profileModel {
	hosts := SSHConfig()
	for _, host := range hosts {
		fmt.Println(host.String())
		fmt.Println("----")
	}
	return profileModel{
		hosts:  hosts,
		cursor: 0,
	}
}

func (m profileModel) Init() tea.Cmd {
	return nil
}

func (m profileModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, cmd
}

func (m profileModel) View() string {
	return "\ntest...\n"
}
