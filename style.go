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
)

var (
	styleConnectionError = lipgloss.NewStyle().
				Foreground(lipgloss.Color("15")).
				Background(lipgloss.Color("1")).
				Bold(true)
	styleTitle = lipgloss.NewStyle().
			Align(lipgloss.Center)
	styleApp = lipgloss.NewStyle().
			Padding(1, 2)
	styleCommand = lipgloss.NewStyle().
			Foreground(lipgloss.Color("0")).
			Background(lipgloss.Color(uranianBlue)).
			Bold(true)
	styleHost = lipgloss.NewStyle().
			Foreground(lipgloss.Color(c500)).
			Bold(true)
	stylePort = styleHost.Copy().
			Foreground(lipgloss.Color(c600))
	styleUser = styleHost.Copy().
			Foreground(lipgloss.Color(c700))
)
