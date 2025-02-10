package ui

import "github.com/charmbracelet/bubbles/key"

// TODO: check if there is a better way to handle the keybindings
// ** There is a lot of repetition here

type irqAddEditKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	goHome key.Binding
	Back   key.Binding
	Select key.Binding
	Help   key.Binding
	Quit   key.Binding
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k irqAddEditKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right, k.Back},
		{k.Select, k.goHome, k.Quit, k.Help},
	}
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k irqAddEditKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.goHome, k.Quit, k.Help}
}

func irqAddEditListKeyMap() *irqAddEditKeyMap {
	return &irqAddEditKeyMap{
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", "left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("→", "right"),
		),
		goHome: key.NewBinding(
			key.WithKeys("home"),
			key.WithHelp("home", "Main menu"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter", "space"),
			key.WithHelp("enter", "select"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
	}
}

type irqKeyMap struct {
	Up         key.Binding
	Down       key.Binding
	Left       key.Binding
	Right      key.Binding
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
			key.WithKeys("up"),
			key.WithHelp("↑", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "move down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", "move left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("→", "move right"),
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
