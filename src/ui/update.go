package ui

import (
	"log"

	"github.com/canonical/rt-conf/src/ui/config"
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

var placeholders_text = []string{
	config.CpuListPlaceholder,
	"Enter on or off",
	config.CpuListPlaceholder,
	config.CpuListPlaceholder,
	config.CpuListPlaceholder,
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Println("\n------------- UPDATE -----------------")
	log.Println("(BEGIN UPDATE) Current menu: ", config.Menu[m.nav.GetCurrMenu()])
	// log.Println("(BEGIN UPDATE) current menu: ", m.currMenu)
	// log.Println("STACK TRACE")
	// log.Println(debug.Stack())
	var cmds []tea.Cmd

	// Menu update handlers
	menuHandlers := map[config.Views]func(tea.KeyMsg) tea.Cmd{
		config.INIT_VIEW_ID:             m.mainMenuUpdate,
		config.KCMD_VIEW_ID:             m.kcmdlineMenuUpdate,
		config.KCMD_CONCLUSSION_VIEW_ID: m.kcmdlineConclussionUpdate,
		config.IRQ_VIEW_ID:              m.IRQMenuUpdate,
		config.IRQ_ADD_EDIT_VIEW_ID:     m.IRQAddEditMenuUpdate,
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// TODO: check if this will not break the app
		h, v := styles.AppStyle.GetFrameSize()
		m.width = msg.Width - h
		m.height = msg.Height - v
		m.main.list.SetSize(msg.Width-h, msg.Height-v)
		m.irq.list.SetSize(msg.Width-h, msg.Height-v)

		m.kcmd.help.Width = msg.Width

	case irqAffinityRule:
		log.Println("(RECEIVED) IRQ Affinity Rule: ", msg)

	case tea.KeyMsg:

		log.Printf("Key pressed: %v", msg.String())
		switch {
		case key.Matches(msg, m.main.keys.Back):
			m.nav.PrevMenu()
			// m.currMenu = m.prevMenu
			return m, nil
			// if !(m.currMenu == mainMenu) {
			// TODO implement prevMenu logic
			// 	m.currMenu = m.prevMenu
			// }
		// This is genertic for all menus

		case key.Matches(msg, m.main.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.main.keys.Home):
			m.nav.ReturnToMainMenu()

		case key.Matches(msg, m.main.keys.Help):
			view := m.nav.GetCurrMenu()
			if view == config.INIT_VIEW_ID {
				// m.main.help.ShowAll = !m.main.help.ShowAll
				m.main.list.SetShowHelp(!m.main.list.ShowHelp())

			} else if view == config.IRQ_VIEW_ID {
				m.irq.list.SetShowHelp(!m.irq.list.ShowHelp())
				// m.kcmd.help.ShowAll = !m.kcmd.help.ShowAll
				// m.irq.list.SetShowHelp(!m.irq.list.ShowHelp())

			} else if view == config.KCMD_VIEW_ID {
				m.kcmd.help.ShowAll = !m.kcmd.help.ShowAll
			} else {
				// m.irq.help.ShowAll = !m.irq.help.ShowAll
			}

		default:
			log.Println("(UPDATE) default handler| Current menu: ",
				config.Menu[m.nav.GetCurrMenu()])

			if handler, exists := menuHandlers[m.nav.GetCurrMenu()]; exists {
				cmds = append(cmds, handler(msg))
			}
		}
	}

	// Update IRQ list view
	new, cmd := m.irq.list.Update(msg)
	m.irq.list = new
	cmds = append(cmds, cmd)

	// Update main menu list view
	newListModel, cmd := m.main.list.Update(msg)
	m.main.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) mainMenuUpdate(msg tea.KeyMsg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.kcmd.Inputs))

	if m.main.list.FilterState() == list.Filtering {
		return tea.Batch(cmds...)
	}
	switch {
	case key.Matches(msg, m.main.keys.Home):
		m.nav.ReturnToMainMenu()

	case key.Matches(msg, m.main.keys.Select):
		selected := m.main.list.SelectedItem().(menuItem)

		// TODO: Improve this selection logic
		// It could be indexed by the menu item
		switch selected.Title() {
		case config.MENU_KCMDLINE:
			log.Println("MENU: Kcmd selected")
			m.nav.SetNewMenu(config.KCMD_VIEW_ID)
			log.Println("Current stack: ", m.nav.PrintMenuStack())

		case config.MENU_IRQAFFINITY:
			log.Println("MENU: IRQ Affinity selected")
			m.nav.SetNewMenu(config.IRQ_VIEW_ID)
			log.Println("Current stack: ", m.nav.PrintMenuStack())
		}

	case key.Matches(msg, m.main.keys.Help):
		m.main.list.SetShowHelp(!m.main.list.ShowHelp())
		// return m, nil
	}
	return tea.Batch(cmds...)
}

func (m Model) IRQMenuUpdate(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	log.Printf("---IRQMenuUpdate: ")
	log.Printf("key pressed: %v", msg.String())
	switch {

	case key.Matches(msg, m.irq.keys.Select):
		log.Printf("confirmed select pressed")
		log.Printf("previous menu: %v", config.Menu[m.nav.GetCurrMenu()])
		m.nav.SetNewMenu(config.IRQ_ADD_EDIT_VIEW_ID)
		log.Println("Current menu: ", config.Menu[m.nav.GetCurrMenu()])
	}
	return cmd
}

func (m Model) IRQAddEditMenuUpdate(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	log.Printf("---IRQAddEditMenuUpdate: ")
	log.Printf("Current menu: %v", config.Menu[m.nav.GetCurrMenu()])
	// TODO: here it should be implemented the logic to add or edit IRQs
	switch {
	default:
		log.Printf("key pressed: %v", msg.String())
	}
	return cmd
}
