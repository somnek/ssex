package main

import (
	"github.com/charmbracelet/bubbles/key"
)

// KeyMap is the mappings of actions to key bindings.
type KeyMap struct {
	Quit       key.Binding
	Next       key.Binding
	Back       key.Binding
    Clear       key.Binding
}

// DefaultKeyMap is the default key map for the application.
var DefaultKeyMap = KeyMap{
	Next:       key.NewBinding(key.WithKeys("enter", "tab"), key.WithHelp("enter", "next")),
	Back:       key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift-tab", "back")),
	Quit:       key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "exit")),
	Clear:       key.NewBinding(key.WithKeys("ctrl+u"), key.WithHelp("ctrl+u", "clear")),
}

// ShortHelp returns a quick help menu.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
        k.Next,
        k.Back,
		k.Quit,
        k.Clear,
	}
}

// FullHelp returns all help options in a more detailed view.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Next, k.Back},
		{k.Quit, k.Clear},
	}
}
