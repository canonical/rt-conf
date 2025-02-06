package ui

import (
	"log"

	"github.com/canonical/rt-conf/src/ui/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
// 	d := list.NewDefaultDelegate()

// 	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
// 		var title string

// 		if i, ok := m.SelectedItem().(item); ok {
// 			title = i.Title()
// 		} else {
// 			return nil
// 		}

// 		switch msg := msg.(type) {
// 		case tea.KeyMsg:
// 			switch {
// 			case key.Matches(msg, keys.choose):
// 				return m.NewStatusMessage(statusMessageStyle("You chose " + title))

// 			case key.Matches(msg, keys.remove):
// 				index := m.Index()
// 				m.RemoveItem(index)
// 				if len(m.Items()) == 0 {
// 					keys.remove.SetEnabled(false)
// 				}
// 				return m.NewStatusMessage(statusMessageStyle("Deleted " + title))
// 			}
// 		}

// 		return nil
// 	}

// 	help := []key.Binding{keys.choose, keys.remove}

// 	d.ShortHelpFunc = func() []key.Binding {
// 		return help
// 	}

// 	d.FullHelpFunc = func() [][]key.Binding {
// 		return [][]key.Binding{help}
// 	}

// 	return d
// }

// // Additional short help entries. This satisfies the help.KeyMap interface and
// // is entirely optional.
// func (d delegateKeyMap) ShortHelp() []key.Binding {
// 	return []key.Binding{
// 		d.choose,
// 		d.remove,
// 	}
// }

// // Additional full help entries. This satisfies the help.KeyMap interface and
// // is entirely optional.
// func (d delegateKeyMap) FullHelp() [][]key.Binding {
// 	return [][]key.Binding{
// 		{
// 			d.choose,
// 			d.remove,
// 		},
// 	}
// }

// // MAIN MODEL

// // TODO: this needs to be replaced by the styles local module
// var (
// 	appStyle = lipgloss.NewStyle().Padding(1, 2)

// 	titleStyle = lipgloss.NewStyle().
// 			Foreground(lipgloss.Color("#FFFDF5")).
// 			Background(lipgloss.Color("#25A065")).
// 			Padding(0, 1)

// 	statusMessageStyle = lipgloss.NewStyle().
// 				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
// 				Render
// )

// type item struct {
// 	title       string
// 	description string
// }

// func (i item) Title() string       { return i.title }
// func (i item) Description() string { return i.description }
// func (i item) FilterValue() string { return i.title }

// func (i *randomItemGenerator) Init() {
// 	i.reset()
// }

// func (r *randomItemGenerator) Size() int {
// 	return len(r.titles)
// }

// type IRQMenuModel struct {
// 	list          list.Model
// 	itemGenerator *randomItemGenerator
// 	keys          *listKeyMap
// 	delegateKeys  *delegateKeyMap
// 	newEntry      bool
// }

// func newModelIRQMenuModel() IRQMenuModel {
// 	var (
// 		itemGenerator randomItemGenerator
// 		delegateKeys  = newDelegateKeyMap()
// 		listKeys      = newListKeyMap()
// 	)

// 	// Make initial list of items
// 	const numItems = 24
// 	items := make([]list.Item, numItems)
// 	for i := 0; i < numItems; i++ {
// 		items[i] = itemGenerator.next()
// 	}

// 	// Setup list
// 	delegate := newItemDelegate(delegateKeys)
// 	irqList := list.New(items, delegate, 0, 0)
// 	irqList.Title = "IRQ Affinity"
// 	irqList.Styles.Title = titleStyle
// 	irqList.AdditionalFullHelpKeys = func() []key.Binding {
// 		return []key.Binding{
// 			listKeys.toggleSpinner,
// 			listKeys.insertItem,
// 			listKeys.toggleTitleBar,
// 			listKeys.toggleStatusBar,
// 			listKeys.togglePagination,
// 			listKeys.toggleHelpMenu,
// 		}
// 	}

// 	return IRQMenuModel{
// 		list:          irqList,
// 		keys:          listKeys,
// 		delegateKeys:  delegateKeys,
// 		itemGenerator: &itemGenerator,
// 	}
// }

// type listKeyMap struct {
// 	toggleSpinner    key.Binding
// 	toggleTitleBar   key.Binding
// 	toggleStatusBar  key.Binding
// 	togglePagination key.Binding
// 	toggleHelpMenu   key.Binding
// 	insertItem       key.Binding
// }

// func newListKeyMap() *listKeyMap {
// 	return &listKeyMap{
// 		insertItem: key.NewBinding(
// 			key.WithKeys("a"),
// 			key.WithHelp("a", "add item"),
// 		),
// 		toggleSpinner: key.NewBinding(
// 			key.WithKeys("s"),
// 			key.WithHelp("s", "toggle spinner"),
// 		),
// 		toggleTitleBar: key.NewBinding(
// 			key.WithKeys("T"),
// 			key.WithHelp("T", "toggle title"),
// 		),
// 		toggleStatusBar: key.NewBinding(
// 			key.WithKeys("S"),
// 			key.WithHelp("S", "toggle status"),
// 		),
// 		togglePagination: key.NewBinding(
// 			key.WithKeys("P"),
// 			key.WithHelp("P", "toggle pagination"),
// 		),
// 		toggleHelpMenu: key.NewBinding(
// 			key.WithKeys("H"),
// 			key.WithHelp("H", "toggle help"),
// 		),
// 	}
// }

// //////////////////////////////////////////////////////////////////////////////

type NewIRQEntryMsg struct {
	irqAffinityRule
}

type irqAffinityRule struct {
	filter, cpulist string
}

func (i irqAffinityRule) Title() string       { return i.filter }
func (i irqAffinityRule) Description() string { return i.cpulist }

// This needs to be implemented to satisfy the list.Item interface
func (i irqAffinityRule) FilterValue() string { return i.filter }

type delegateKeyMap struct {
	choose key.Binding
	remove key.Binding
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
		remove: key.NewBinding(
			key.WithKeys("x", "backspace"),
			key.WithHelp("x", "delete"),
		),
	}
}

func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	log.Println("(newItemDelegate) Creating new item delegate")
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var title string

		if i, ok := m.SelectedItem().(irqAffinityRule); ok {
			title = i.Title()
		} else {
			return nil
		}
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.choose):
				return m.NewStatusMessage(
					styles.StatusMessageStyle("You chose " + title))

			case key.Matches(msg, keys.remove):
				index := m.Index()
				m.RemoveItem(index)
				if len(m.Items()) == 0 {
					keys.remove.SetEnabled(false)
				}
				return m.NewStatusMessage(
					styles.StatusMessageStyle("Deleted " + title))
			}
		}
		return nil
	}

	help := []key.Binding{keys.choose, keys.remove}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}
	return d
}
