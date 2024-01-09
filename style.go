package main

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	// Di Sierra
	c50         = "#FBF7F1"
	c100        = "#F6ECDE"
	c200        = "#ECD6BC"
	c300        = "#E0B991"
	c400        = "#CF8C56"
	c500        = "#C97B46"
	c600        = "#BB663B"
	c700        = "#9C5132"
	c800        = "#7E422E"
	c900        = "#663828"
	c950        = "#361B14"
	onyx        = "#36393B"
	uranianBlue = "#A5D8FF"
	cadetGray   = "#8DA7BE"
	slateGray   = "#717C89"
	cordovan    = "#90323D"
	clared      = "#840032"
	thulianPink = "#DE639A"

	// occupied by non empty lines
	reservedLinesHeight = 11
)

var (
	styleTitle = lipgloss.NewStyle().
			Align(lipgloss.Center)
	styleApp = lipgloss.NewStyle().
			Padding(0, 1, 0, 1)
	styleCommand = lipgloss.NewStyle().
			Foreground(lipgloss.Color("232")).
			Background(lipgloss.Color(uranianBlue)).
			Bold(true)
	styleError = lipgloss.NewStyle().
			Foreground(lipgloss.Color("232")).
			Background(lipgloss.Color(thulianPink)).
			Bold(true)
	styleHost = lipgloss.NewStyle().
			Foreground(lipgloss.Color(c400)).
			Bold(true)
	stylePort = styleHost.Copy().
			Foreground(lipgloss.Color(c500))
	styleUser = styleHost.Copy().
			Foreground(lipgloss.Color(c600))
	styleConnectedStr = lipgloss.NewStyle().
				Bold(true)
	styleWordConnected = lipgloss.NewStyle().
				Foreground(lipgloss.Color(c200))
	styleCheckMark = lipgloss.NewStyle().
			Foreground(lipgloss.Color(uranianBlue)).
			Bold(true)
	styleConnectionError = lipgloss.NewStyle().
				Foreground(lipgloss.Color(thulianPink)).
				Bold(true)
)
