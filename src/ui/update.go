package ui

import (
	"github.com/canonical/rt-conf/src/ui/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// TODO: The values of kcmdline menu should come from the YAML file

// TODO: Improve navigation on kcmdline menu.
// * NOTE: For now, the kcmdline menu navigation is shared between components
// * that are part of a multiple text input form and two buttons.
// * The indexing is done by the focusIndex variable, which is incremented or
// * decremented by the NextIndex() and PrevIndex() functions.
// * But once it's needed to apply functions specific  to the text inputs, the
// * it's necessary to check everytime for overflow and underflow of the
// * focusIndex

const (
	isolcpusIndex = iota
	nohzIndex
	nohzFullIndex
	kthreadsCPUsIndex
	irqaffinityIndex
	applyButtonIndex
	backButtonIndex
)

const cpuListPlaceholder = "Enter a CPU list like: 4-n or 3-5 or 2,4,5 "
const irqFilterPlaceholder = "Insert filter parameters for IRQs"

var placeholders_text = []string{
	cpuListPlaceholder,
	"Enter on or off",
	cpuListPlaceholder,
	cpuListPlaceholder,
	cpuListPlaceholder,
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Menu update handlers
	menuHandlers := map[menuOpt]func(tea.KeyMsg) tea.Cmd{
		mainMenu:        m.updateMainMenu,
		kcmdlineMenu:    m.updateKcmdlineMenu,
		irqAffinityMenu: m.updateIRQMenu,
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// TODO: check if this will not break the app
		h, v := styles.AppStyle.GetFrameSize()

		m.width = msg.Width - h
		m.height = msg.Height - v

		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.help.Width = msg.Width

	case tea.KeyMsg:
		switch {
		// This is genertic for all menus
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.goHome):
			m.currMenu = mainMenu
		case key.Matches(msg, m.keys.Help):
			if m.currMenu == mainMenu {
				m.list.SetShowHelp(!m.list.ShowHelp())
			} else {
				m.help.ShowAll = !m.help.ShowAll
			}
		default:
			if handler, exists := menuHandlers[m.currMenu]; exists {
				cmds = append(cmds, handler(msg))
			}
		}
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *Model) updateMainMenu(msg tea.KeyMsg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.kcmdInputs))

	if m.list.FilterState() == list.Filtering {
		return tea.Batch(cmds...)
	}
	switch {
	case key.Matches(msg, m.keys.goHome):
		m.currMenu = mainMenu

	case key.Matches(msg, m.keys.Select):
		selected := m.list.SelectedItem().(item)

		// TODO: Improve this selection logic
		// It could be indexed by the menu item
		switch selected.Title() {
		case "Kernel cmdline":
			m.prevMenu = mainMenu
			m.currMenu = kcmdlineMenu
		case "IRQ Affinity":
			m.prevMenu = mainMenu
			m.currMenu = irqAffinityMenu
		}

	case key.Matches(msg, m.keys.Help):
		m.list.SetShowHelp(!m.list.ShowHelp())
		// return m, nil
	}
	return tea.Batch(cmds...)
}
