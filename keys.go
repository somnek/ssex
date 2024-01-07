package main

import (
	"github.com/charmbracelet/bubbles/key"
)

// KeyMap is the mappings of actions to key bindings.
type KeyMap struct {
	Quit       key.Binding
	Next       key.Binding
	Back       key.Binding
	ToggleHelp key.Binding
}

// DefaultKeyMap is the default key map for the application.
var DefaultKeyMap = KeyMap{
	Quit:       key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "exit")),
	Next:       key.NewBinding(key.WithKeys("enter", "tab"), key.WithHelp("enter", "next")),
	Back:       key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift-tab", "back")),
	ToggleHelp: key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
}

// ShortHelp returns a quick help menu.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Quit,
		k.ToggleHelp,
	}
}

// FullHelp returns all help options in a more detailed view.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.ToggleHelp},
		{k.Next, k.Back},
	}
}
