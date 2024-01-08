package main

import (
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/quick"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sshModel struct {
	client  *Client
	profile Profile
	command string // executed command
	input   textinput.Model
	output  string
	height  int
	width   int
	help    help.Model
	keys    KeyMap
	err     error
}

type errMsg error

func initSSHModel(client *Client, p Profile, height, width int) sshModel {

	// input
	t := textinput.New()
	t.Placeholder = "Enter command"
	t.Focus()

	keys := DefaultKeyMap
	keys.Next.Unbind()
	keys.Back.Unbind()

	return sshModel{
		profile: p,
		keys:    keys,
		help:    help.New(),
		client:  client,
		height:  height,
		width:   width,
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
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			input := m.input.Value()
			output, err := m.client.RunCmd(input)

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

	// title
	title := buildTitle(m.width)
	b.WriteString(title)
	b.WriteString("\n")

	// connection
	connStr := buildConnStr(m.profile, m.width)
	b.WriteString(connStr)
	b.WriteString("\n")
	// paddingLen := horizontalPadLength(connStr, m.width)
	// styleConnStr.PaddingLeft(paddingLen)
	// b.WriteString(styleConnStr.Render(connStr) + "\n")

	// input
	b.WriteString(m.input.View())
	b.WriteString("\n\n")

	// command
	if m.command != "" {
		commandLeft := styleCommand.Render("Command:")
		b.WriteString(commandLeft + " " + m.command + "\n\n")
	}

	// output
	if m.output != "" {
		b.WriteString(m.output)
		b.WriteString("\n\n")
	}

	// pad with empty lines
	buildEmptyLine(&b, m.height)

	// help
	b.WriteString(m.help.View(m.keys))
	return styleApp.Render(b.String())
}

func buildOutput(output string, height int) string {
	ob := strings.Builder{}
	quick.Highlight(&ob, output, "actionscript 3", "terminal16m", "friendly")

	// trim long output
	lines := strings.Split(ob.String(), "\n")
	if len(lines) > height-10 {
		truncText := fmt.Sprintf("... %d more lines truncated ...", len(lines)-(height+10))
		lines = lines[:height-10]
		return strings.Join(lines, "\n") + "\n" + truncText
	}

	return ob.String()
}

func buildConnStr(p Profile, width int) string {
	host := p.Host
	port := p.Port
	user := p.User
	hostRender := styleHost.Render(host)
	portRender := stylePort.Render(port)
	userRender := styleUser.Render(user)
	connectedStrRender := lipgloss.NewStyle().Foreground(lipgloss.Color(c200)).Render("Connected")

	connStr := fmt.Sprintf(
		"%s to %s:%s as %s",
		connectedStrRender,
		hostRender,
		portRender,
		userRender,
	)
	if width > 0 {
		paddingLen := horizontalPadLength(connStr, width)
		styleConnStr.PaddingLeft(paddingLen)
	}
	return styleConnStr.Render(connStr)
}
