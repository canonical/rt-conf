package ui

import "github.com/charmbracelet/bubbles/key"

type irqKeyMap struct {
	Up         key.Binding
	Down       key.Binding
	goHome     key.Binding
	Select     key.Binding
	Help       key.Binding
	Quit       key.Binding
	CursorMode key.Binding
	Back       key.Binding
	Add        key.Binding
	Remove     key.Binding
	Apply      key.Binding
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k irqKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back, k.Apply, k.Add, k.Remove, k.Help}, // first column
		{k.Up, k.Down, k.goHome, k.Back, k.Quit},   // second column
	}
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k irqKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.Apply, k.Add, k.Remove, k.Help}
}

func irqMenuListKeyMap() *irqKeyMap {
	return &irqKeyMap{
		Add: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add item"),
		),
		Remove: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp("x", "remove item"),
		),
		Apply: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "add item"),
		),
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
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q/ctrl+c", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
	}
}
