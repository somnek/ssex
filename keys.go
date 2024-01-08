package main

import (
	"github.com/charmbracelet/bubbles/key"
)

// KeyMap is the mappings of actions to key bindings.
type KeyMap struct {
	Quit   key.Binding
	Next   key.Binding
	Back   key.Binding
	Clear  key.Binding
	Cancel key.Binding
	Enter  key.Binding
}

// TODO: separate out the keymap for the form and the panel
// DefaultKeyMap is the default key map for the application.
var DefaultKeyMap = KeyMap{
	Next:   key.NewBinding(key.WithKeys("enter", "tab"), key.WithHelp("enter", "next")),
	Back:   key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift-tab", "back")),
	Quit:   key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "exit")),
	Clear:  key.NewBinding(key.WithKeys("ctrl+u"), key.WithHelp("ctrl+u", "clear")),
	Cancel: key.NewBinding(key.WithKeys("esc", "q"), key.WithHelp("esc", "cancel/new connection"), key.WithDisabled()),
	Enter:  key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit"), key.WithDisabled()),
}

// ShortHelp returns a quick help menu.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Next,
		k.Back,
		k.Quit,
		k.Clear,
		k.Cancel,
		k.Enter,
	}
}

// FullHelp returns all help options in a more detailed view.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Next, k.Back},
		{k.Quit, k.Clear},
		{k.Cancel, k.Enter},
	}
}
