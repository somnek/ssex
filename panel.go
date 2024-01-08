package main

import (
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/quick"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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
	// t.Validate = inputValidator
	t.Focus()

	keys := DefaultKeyMap
	keys.Next.Unbind()
	keys.Back.Unbind()
	keys.Enter.SetEnabled(true)

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
	var err error

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Enter):
			input := m.input.Value()
			err = inputValidator(input)
			if err != nil {
				m.err = err
				return m, nil
			}

			output, _ := m.client.RunCmd(input)
			m.output = buildOutput(output, m.height)
			m.command = input
			m.input.SetValue("")
			m.err = nil

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

	// input
	b.WriteString(m.input.View())
	b.WriteString("\n\n")

	if m.err != nil {
		// errors
		errorLeft := styleError.Render("Error:")
		b.WriteString(errorLeft + " " + m.err.Error() + "\n")
	} else {
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
	empties := height - reservedLinesHeight

	if len(lines) > empties {
		truncText := fmt.Sprintf("... %d more lines truncated ...", len(lines)-(height+10))
		lines = lines[:empties]
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
	wordConnectedRender := styleWordConnected.Render("Connected")
	checkMarkRender := styleCheckMark.Render("✔")

	connStr := fmt.Sprintf(
		"%s %s to %s:%s as %s",
		checkMarkRender,
		wordConnectedRender,
		hostRender,
		portRender,
		userRender,
	)
	if width > 0 {
		paddingLen := horizontalPadLength(connStr, width)
		styleConnectedStr.PaddingLeft(paddingLen)
	}
	return styleConnectedStr.Render(connStr)
}

func inputValidator(s string) error {
	if s == "" {
		return errMsg(fmt.Errorf("command cannot be empty"))
	}
	return nil
}
