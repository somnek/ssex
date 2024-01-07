package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/alecthomas/chroma/quick"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type sshModel struct {
	client  *Client
	command string // executed command
	input   textinput.Model
	output  string
	height  int
	width  int
	help    help.Model
	keys    KeyMap
	err     error
}

type errMsg error

func initSSHModel(client *Client, height, width int) sshModel {

	// input
	t := textinput.New()
	t.Placeholder = "Enter command"
	t.Focus()

	return sshModel{
		keys:    DefaultKeyMap,
		help:    help.New(),
		client:  client,
		height:  height,
		width:  width,
		input:   t,
	}
}

func (m sshModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m sshModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
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

			// m.output = output
            m.output = buildOutput(output, m.height)
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
	return m, cmd
}

func (m sshModel) View() string {

	var b strings.Builder
	title := buildTitle(m.width)
	b.WriteString(title)

	b.WriteString(m.input.View())
	b.WriteString("\n\n")

	if m.command != "" {
		commandLeft := styleCommand.Render("Command:")
		b.WriteString(commandLeft + " " + m.command + "\n\n")
	}
	// new writer
    if m.output != ""{
        // output := buildOutput(m.output, m.height)
        // b.WriteString(output)
        b.WriteString(m.output)
        b.WriteString("\n\n")
    }

	buildEmptyLine(&b, m.height)
	b.WriteString(m.help.View(m.keys))
	return styleApp.Render(b.String())
}

func buildOutput(output string, height int) string{
	ob := strings.Builder{}
	quick.Highlight(&ob, output, "actionscript 3", "terminal16m", "friendly")

    // trim long output
    lines := strings.Split(ob.String(), "\n")
    if len(lines) > height - 10{
        truncText := fmt.Sprintf("... %d more lines truncated ...", len(lines) - (height+10))
        lines = lines[:height-10]
        return strings.Join(lines, "\n") + "\n" + truncText
    }

    return ob.String()
}
