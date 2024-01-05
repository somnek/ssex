package main

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

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
}

func initialModel() profileModel {
	var p Profile

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
				Placeholder("e.g 22 (default)").
				Validate(func(s string) error {
					if s == "" {
						return errors.New("port is required")
					}
					return nil
				}).
				Description("Specify the port number for the remote server:"),
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
				Description("Provide the user for accessing the remote server:"),
		),
	)

	return profileModel{
		form:    f,
		profile: p,
	}

}

func (m profileModel) Init() tea.Cmd {
	return m.form.Init()
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
		case "X":
			return initSSHModel(), nil
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, cmd
}

func (m profileModel) View() string {
	// move this part to Update to check if form is completed
	if m.form.State == huh.StateCompleted {
		host := m.form.GetString("Host")
		port := m.form.GetString("Port")
		user := m.form.GetString("User")
		return "\n" + host + "\n" + port + "\n" + user
	}
	return m.form.View()
}

// parsedConfig := ParseSSHConfig()
// var hosts []Host
// for _, block := range parsedConfig {
// 	hosts = append(hosts, Host{
// 		Alias: block,
// 	})
// 	fmt.Println(block.Nodes)
// 	fmt.Println("----")
// }
