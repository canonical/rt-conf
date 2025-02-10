package ui

import (
	cmp "github.com/canonical/rt-conf/src/ui/components"
	"github.com/canonical/rt-conf/src/ui/config"
	"github.com/canonical/rt-conf/src/ui/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type MainMenuModel struct {
	Nav          *cmp.MenuNav
	Width        int
	Height       int
	keys         *KeyMap
	list         list.Model
	delegateKeys *selectKeyMap
}

type KeyMap struct {
	Home   key.Binding
	Back   key.Binding
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Help   key.Binding
	Quit   key.Binding
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select, k.Quit, k.Help}, // first column
	}
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Help}
}

func newMainMenuListKeyMap() *KeyMap {
	return &KeyMap{
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter", "space"),
			key.WithHelp("enter", "select menu"),
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

type menuItem struct {
	title       string
	description string
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.description }
func (i menuItem) FilterValue() string { return i.title }
func (r *menuItem) Size() int          { return config.NUMBER_OF_MENUS }
func (r *menuItem) Init()              {}

func NewMainMenuModel() MainMenuModel {
	keys := newMainMenuListKeyMap()
	var delegateKeys = newDelegateKeyMapMainMenu()

	items := []list.Item{
		menuItem{config.MENU_KCMDLINE, config.DESC_KCMDLINE},
		menuItem{config.MENU_IRQAFFINITY, config.DESC_IRQAFFINITY},
		menuItem{config.MENU_PWRMGMT, config.DESC_KCMDLINE},
	}

	delegate := newItemDelegateMainMenu(delegateKeys)
	delegate.Styles.SelectedDesc = styles.SelectedDesc
	delegate.Styles.SelectedTitle = styles.SelectedTitle

	delegate.Styles.NormalDesc = styles.NormalDesc
	delegate.Styles.NormalTitle = styles.NormalTitle

	delegate.Styles.DimmedDesc = styles.DimmedDesc
	delegate.Styles.DimmedTitle = styles.DimmedTitle

	delegate.Styles.FilterMatch = styles.FilterMatch
	menuList := list.New(items, delegate, 0, 0)
	menuList.SetShowHelp(true)
	menuList.Title = "rt-conf tool"
	menuList.Styles.Title = styles.TitleStyle
	menuList.SetShowStatusBar(false)
	menuList.KeyMap.NextPage.SetEnabled(false)
	menuList.KeyMap.PrevPage.SetEnabled(false)
	menuList.KeyMap.Filter.SetEnabled(false)

	nav := cmp.GetMenuNavInstance()
	return MainMenuModel{
		Nav:          nav,
		keys:         keys,
		list:         menuList,
		delegateKeys: delegateKeys,
	}
}

func (m MainMenuModel) Init() tea.Cmd { return nil }

func (m MainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 5)
	// log.Println("---MainMenuUpdate: ")

	if m.list.FilterState() == list.Filtering {
		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {

		case key.Matches(msg, m.keys.Home):
			m.Nav.ReturnToMainMenu()

		case key.Matches(msg, m.keys.Select):
			selected := m.list.SelectedItem().(menuItem)

			// TODO: Improve this selection logic
			// It could be indexed by the menu item
			switch selected.Title() {
			case config.MENU_KCMDLINE:
				// log.Println("MENU: Kcmd selected")
				m.Nav.SetNewMenu(config.KCMD_VIEW_ID)
				// log.Println("Current stack: ", m.Nav.PrintMenuStack())

			case config.MENU_IRQAFFINITY:
				// log.Println("MENU: IRQ Affinity selected")
				m.Nav.SetNewMenu(config.IRQ_VIEW_ID)
				// log.Println("Current stack: ", m.Nav.PrintMenuStack())
			}

		case key.Matches(msg, m.keys.Help):
			m.list.SetShowHelp(!m.list.ShowHelp())
			// return m, nil
		}
	}

	return m, tea.Batch(cmds...)
}

func (m MainMenuModel) View() string {
	return styles.AppStyle.Render(m.list.View())
}
