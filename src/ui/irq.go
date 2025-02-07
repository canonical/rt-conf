package ui

import (
	"log"

	"github.com/canonical/rt-conf/src/ui/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// TODO: add conclussion screen saying:
// "5 IRQ filter rules are aplied to the system"
// [ BACK ]

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
