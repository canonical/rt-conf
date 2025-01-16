package ui

import "github.com/charmbracelet/bubbles/key"

type listKeyMap struct {
	Up         key.Binding
	Down       key.Binding
	goHome     key.Binding
	Select     key.Binding
	Help       key.Binding
	Quit       key.Binding
	CursorMode key.Binding
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k listKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.goHome},   // first column
		{k.Help, k.Select, k.Quit}, // second column
	}
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k listKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{

		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		goHome: key.NewBinding(
			key.WithKeys("g", "home"),
			key.WithHelp("home/g", "Main menu"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter", "space"),
			key.WithHelp("enter", "select menu"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctl+c"),
			key.WithHelp("q/ctrl+c", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
	}
}
