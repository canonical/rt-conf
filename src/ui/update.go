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

func (m IRQAddEditMenu) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Println("\n------------- UPDATE -----------------")
	log.Println("(BEGIN UPDATE) Current menu: ", config.Menu[m.Nav.GetCurrMenu()])
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		// TODO: check if this will not break the app
		h, v := styles.AppStyle.GetFrameSize()
		m.width = msg.Width - h
		m.height = msg.Height - v

		// m.kcmd.list.SetSize(msg.Width-h, msg.Height-v) // TODO: implement this
		m.main.list.SetSize(msg.Width-h, msg.Height-v)
		m.irq.list.SetSize(msg.Width-h, msg.Height-v)
		m.kcmd.help.Width = msg.Width

	case irqAffinityRule:
		log.Println("(RECEIVED) IRQ Affinity Rule: ", msg)

	case tea.KeyMsg:
		log.Printf("Key pressed: %v", msg.String())
		switch {
		case key.Matches(msg, m.main.keys.Back):
			m.Nav.PrevMenu()
			return m, nil
			// if !(m.currMenu == mainMenu) {
			// TODO implement prevMenu logic
			// 	m.currMenu = m.prevMenu
			// }
		// This is genertic for all menus

		case key.Matches(msg, m.main.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.main.keys.Home):
			m.Nav.ReturnToMainMenu()

		case key.Matches(msg, m.main.keys.Help):
			view := m.Nav.GetCurrMenu()
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

			// default:
			// 	log.Println("(UPDATE) default handler| Current menu: ",
			// 		config.Menu[m.Nav.GetCurrMenu()])

			// if handler, exists := menuHandlers[m.Nav.GetCurrMenu()]; exists {
			// 	cmds = append(cmds, handler(msg))
			// }
		}
	}

	activeMenu := m.GetActiveMenu()
	var cmd tea.Cmd
	// Calling the Update() function of the active menu
	m.currMenu, cmd = activeMenu.Update(msg)
	cmds = append(cmds, cmd)

	// Update main menu list view
	log.Println("-> Updating main menu list:")
	newListModel, cmd := m.main.list.Update(msg)
	m.main.list = newListModel
	cmds = append(cmds, cmd)

	// Update IRQ list view
	new, cmd := m.irq.list.Update(msg)
	m.irq.list = new
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *MainMenuModel) Update(msg tea.Msg) (Menu, tea.Cmd) {
	cmds := make([]tea.Cmd, 5)
	log.Println("---MainMenuUpdate: ")

	if m.list.FilterState() == list.Filtering {
		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := styles.AppStyle.GetFrameSize()
		m.Width = msg.Width - h
		m.Height = msg.Height - v
		m.list.SetSize(msg.Width-h, msg.Height-v)

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
				log.Println("MENU: Kcmd selected")
				m.Nav.SetNewMenu(config.KCMD_VIEW_ID)
				log.Println("Current stack: ", m.Nav.PrintMenuStack())
				// return KcmdlineMenuModel{}, nil

			case config.MENU_IRQAFFINITY:
				log.Println("MENU: IRQ Affinity selected")
				m.Nav.SetNewMenu(config.IRQ_VIEW_ID)
				log.Println("Current stack: ", m.Nav.PrintMenuStack())
			}

		case key.Matches(msg, m.keys.Help):
			m.list.SetShowHelp(!m.list.ShowHelp())
			// return m, nil
		}
	}

	return m, tea.Batch(cmds...)
}

func (m IRQMenuModel) Update(msg tea.Msg) (Menu, tea.Cmd) {
	var cmd tea.Cmd
	log.Printf("---IRQMenuUpdate: ")

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := styles.AppStyle.GetFrameSize()
		m.width = msg.Width - h
		m.height = msg.Height - v

		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		log.Printf("key pressed: %v", msg.String())

		switch {

		case key.Matches(msg, m.irq.keys.Select):
			log.Printf("confirmed select pressed")
			log.Printf("previous menu: %v", config.Menu[m.Nav.GetCurrMenu()])
			m.Nav.SetNewMenu(config.IRQ_ADD_EDIT_VIEW_ID)
			log.Println("Current menu: ", config.Menu[m.Nav.GetCurrMenu()])
		}
	}
	return m, cmd
}

func (m IRQAddEditMenu) Update(msg tea.Msg) (Menu, tea.Cmd) {
	var cmd tea.Cmd
	log.Printf("---IRQAddEditMenuUpdate: ")
	log.Printf("Current menu: %v", config.Menu[m.Nav.GetCurrMenu()])
	// TODO: here it should be implemented the logic to add or edit IRQs
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := styles.AppStyle.GetFrameSize()
		m.Width = msg.Width - h
		m.Height = msg.Height - v

		// m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		switch {
		default:
			log.Printf("key pressed: %v", msg.String())
		}
	}
	return m, cmd
}
