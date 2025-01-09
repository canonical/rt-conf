package ui

import "github.com/charmbracelet/bubbles/key"

type listKeyMap struct {
	Up         key.Binding
	Down       key.Binding
	goHome     key.Binding
	selectMenu key.Binding
	Help       key.Binding
	Quit       key.Binding
}

func (k listKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down}, // first column
		{k.Help},       // second column
	}
}

func (k listKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		// insertItem: key.NewBinding(
		// 	key.WithKeys("a"),
		// 	key.WithHelp("a", "add item"),
		// ),
		// toggleTitleBar: key.NewBinding(
		// 	key.WithKeys("T"),
		// 	key.WithHelp("T", "toggle title"),
		// ),
		goHome: key.NewBinding(
			key.WithKeys("g", "home"),
			key.WithHelp("g", "Home screen"),
		),
		selectMenu: key.NewBinding(
			key.WithKeys("enter", "space"),
			key.WithHelp("enter", "select menu"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctl+c"),
			key.WithHelp("q", "quit"),
		),
		// toggleSpinner: key.NewBinding(
		// 	key.WithKeys("s"),
		// 	key.WithHelp("s", "toggle spinner"),
		// ),
		// togglePagination: key.NewBinding(
		// 	key.WithKeys("P"),
		// 	key.WithHelp("P", "toggle pagination"),
		// ),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
	}
}
