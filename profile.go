package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type profileModel struct {
	hosts  []Host
	cursor int
	err    error
}

type Host struct {
	Alias        string
	HostName     string
	User         string
	Port         string
	IdentityFile string
}

func initialModel() profileModel {
	parsedConfig := ParseSSHConfig()
	// var hosts []Host
	for _, block := range parsedConfig {
		// hosts = append(hosts, Host{
		// 	Alias: block
		// })
		fmt.Println(block.Nodes)
		fmt.Println("----")
	}
	return profileModel{
		hosts:  []Host{},
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
