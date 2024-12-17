package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	strongOrange = "#E95420"
)

var (
	appStyle = lipgloss.
			NewStyle().
			Padding(1, 2).
			Foreground(lipgloss.Color("#2D3748"))

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color(strongOrange)).
			Padding(0, 1)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFDF5")).
				Background(lipgloss.Color(strongOrange)).
				Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(
			lipgloss.AdaptiveColor{
				Light: "#3EB34F",
				Dark:  "#3EB34F",
			},
		).
		Render
)

type item struct {
	title       string
	description string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type listKeyMap struct {
	// toggleTitleBar   key.Binding
	// insertItem       key.Binding
	toggleSpinner    key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
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
		toggleSpinner: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "toggle spinner"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}

type Model struct {
	list          list.Model
	itemGenerator *menuItems
	keys          *listKeyMap
	delegateKeys  *selectKeyMap
}

func NewModel() Model {
	var (
		itemGenerator menuItems
		delegateKeys  = newDelegateKeyMap()
		listKeys      = newListKeyMap()
	)

	// Make initial list of items
	const numItems = 2
	items := make([]list.Item, numItems)
	for i := 0; i < numItems; i++ {
		items[i] = itemGenerator.next()
	}

	// Setup list
	delegate := newItemDelegate(delegateKeys)
	menuList := list.New(items, delegate, 0, 0)
	menuList.Title = "rt-conf tool"
	menuList.Styles.Title = titleStyle

	menuList.SetShowStatusBar(false)
	menuList.Styles.PaginationStyle = lipgloss.NewStyle().Background(lipgloss.Color("#CDCDCD"))
	menuList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.toggleSpinner,
			// listKeys.insertItem,
			// listKeys.toggleTitleBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}

	return Model{
		list:          menuList,
		keys:          listKeys,
		delegateKeys:  delegateKeys,
		itemGenerator: &itemGenerator,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.toggleSpinner):
			cmd := m.list.ToggleSpinner()
			return m, cmd

		// case key.Matches(msg, m.keys.toggleTitleBar):
		// 	v := !m.list.ShowTitle()
		// 	m.list.SetShowTitle(v)
		// 	m.list.SetShowFilter(v)
		// 	m.list.SetFilteringEnabled(v)
		// 	return m, nil

		case key.Matches(msg, m.keys.togglePagination):
			m.list.SetShowPagination(!m.list.ShowPagination())
			return m, nil

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil

			// case key.Matches(msg, m.keys.insertItem):
			// 	m.delegateKeys.choose.SetEnabled(true)
			// 	newItem := m.itemGenerator.next()
			// 	insCmd := m.list.InsertItem(0, newItem)
			// 	statusCmd := m.list.NewStatusMessage(statusMessageStyle("Added " + newItem.Title()))
			// 	return m, tea.Batch(insCmd, statusCmd)
		}
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return appStyle.Render(m.list.View())
}
