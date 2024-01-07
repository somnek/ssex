package main

import (
	"log"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sshModel struct {
	client  *Client
	command string // executed command
	input   textinput.Model
	output  string
	spinner spinner.Model
	height  int
	help    help.Model
	keys    KeyMap
	err     error
}

type errMsg error

func initSSHModel(client *Client, height int) sshModel {

	// input
	t := textinput.New()
	t.Placeholder = "Enter command"
	t.Focus()

	// spinner
	s := spinner.New()
	s.Spinner = spinner.Pulse
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(cadetGray))

	return sshModel{
		keys:    DefaultKeyMap,
		help:    help.New(),
		client:  client,
		height:  height,
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
	case tea.WindowSizeMsg:
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		// switch {
		// case key.Matches(msg, m.keys.ToggleHelp):
		// 	m.help.ShowAll = !m.help.ShowAll
		// }
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
			m.command = input
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
	title := buildTitle()
	b.WriteString(title)

	b.WriteString(m.input.View())
	b.WriteString("\n\n")

	if m.command != "" {
		commandLeft := styleCommand.Render("Command:")
		b.WriteString(commandLeft + " " + m.command + "\n\n")
	}
	// new writer
	ob := strings.Builder{}
	quick.Highlight(&ob, m.output, "actionscript 3", "terminal16m", "friendly")
	b.WriteString(ob.String())
	b.WriteString("\n\n")

	b.WriteString(m.spinner.View())
	buildEmptyLine(&b, m.height)
	b.WriteString(m.help.View(m.keys))
	return styleApp.Render(b.String())
}
