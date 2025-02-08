package ui

import (
	"log"

	cmp "github.com/canonical/rt-conf/src/ui/components"
	"github.com/canonical/rt-conf/src/ui/config"
	"github.com/canonical/rt-conf/src/ui/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var placeholders_text = []string{
	config.CpuListPlaceholder,
	"Enter on or off",
	config.CpuListPlaceholder,
	config.CpuListPlaceholder,
	config.CpuListPlaceholder,
}

const (
	irqFilterIndex = iota
	cpuListIndex
	applyBtnIndex
	backBtnIndex
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Println("\n------------- UPDATE -----------------")
	log.Println("(BEGIN UPDATE) Current menu: ", config.Menu[m.Nav.GetCurrMenu()])
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		// ** NOTE: This needs to be done here in the main Update() function
		// ** because for some reason if passed for the other
		// ** Update() functions the size of the window is not updated
		// TODO: check why this is happening
		h, v := styles.AppStyle.GetFrameSize()
		m.Width = msg.Width - h
		m.Height = msg.Height - v
		// m.kcmd.list.SetSize(msg.Width-h, msg.Height-v) // TODO: implement this
		m.main.list.SetSize(msg.Width-h, msg.Height-v)
		m.irq.list.SetSize(msg.Width-h, msg.Height-v)
		m.kcmd.help.Width = msg.Width
		m.kcmd.Width = msg.Width - h
		m.kcmd.Height = msg.Height - v
		m.irq.Width = msg.Width - h
		m.irq.Height = msg.Height - v
		m.irq.irq.width = msg.Width - h
		m.irq.irq.height = msg.Height - v

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

	// Update Kcmdline list view
	// newKcmd, cmd := m.kcmd.list.Update(msg)
	// m.kcmd.list = newKcmd
	// cmds = append(cmds, cmd)

	// Update IRQ list view
	new, cmd := m.irq.list.Update(msg)
	m.irq.list = new
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m MainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m IRQMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	log.Printf("---IRQMenuUpdate: ")

	switch msg := msg.(type) {
	// case tea.WindowSizeMsg:
	// 	h, v := styles.AppStyle.GetFrameSize()
	// 	m.Width = msg.Width - h
	// 	m.Height = msg.Height - v
	// 	m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		log.Printf("key pressed: %v", msg.String())

		switch {
		case key.Matches(msg, m.keys.Apply):
			log.Println("-> Entry: Apply changes")
			num := len(m.list.Items())
			if num >= 1 {
				m.conclussion.num = num
				m.conclussion.logMsg =
					"IRQ affinity rules are aplied to the system"
				m.Nav.SetNewMenu(config.IRQ_CONCLUSSION_VIEW_ID)
			}

		case key.Matches(msg, m.irq.keys.Select):
			log.Printf("confirmed select pressed")
			log.Printf("previous menu: %v", config.Menu[m.Nav.GetCurrMenu()])
			m.Nav.SetNewMenu(config.IRQ_ADD_EDIT_VIEW_ID)
			log.Println("Current menu: ", config.Menu[m.Nav.GetCurrMenu()])
		}
	}
	return m, cmd
}

func (m IRQConclussion) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Select):
			m.Nav.PrevMenu()
		}
	}
	return m, nil
}

func (m *IRQAddEditMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Println("---(IRQAddEditMenu - start")
	cmds := make([]tea.Cmd, len(m.Inputs))

	totalIrqItems := len(m.Inputs) + 2
	index := cmp.NewNavigation(&m.FocusIndex, &totalIrqItems)
	log.Println("\tFocus index: ", m.FocusIndex)
	log.Println("\tAddres: ", &m.FocusIndex)

	log.Printf("---IRQAddEditMenuUpdate: ")
	log.Printf("Current menu: %v", config.Menu[m.Nav.GetCurrMenu()])
	// TODO: here it should be implemented the logic to add or edit IRQs
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Select),
			key.Matches(msg, m.keys.Up),
			key.Matches(msg, m.keys.Down),
			key.Matches(msg, m.keys.Left),
			key.Matches(msg, m.keys.Right):

			// Handle navigation between the buttons
			if m.FocusIndex == applyBtnIndex &&
				key.Matches(msg, m.keys.Left) {
				index.Prev()
			}
			if m.FocusIndex == applyBtnIndex &&
				key.Matches(msg, m.keys.Right) {
				index.Next()
			}
			if m.FocusIndex == backBtnIndex &&
				key.Matches(msg, m.keys.Right) {
				index.Next()
			}
			if m.FocusIndex == backBtnIndex &&
				key.Matches(msg, m.keys.Left) {
				index.Prev()
			}

			// Handle [ Back ] button
			if m.FocusIndex == backBtnIndex &&
				key.Matches(msg, m.keys.Select) {
				// log.Println("pressed [ BACK ]: Back to main menu")
				m.Nav.PrevMenu()
			}

			// Did the user press enter while the apply button was focused?
			// TODO: improve mapping of len(m.inputs) to the apply button
			if key.Matches(msg, m.keys.Select) &&
				m.FocusIndex == len(m.Inputs) {

				var empty int
				for i := range m.Inputs {
					v := m.Inputs[i].Value()
					if v == "" {
						empty++
					}
				}

				if empty == len(m.Inputs) {
					m.errorMsgFilter = "\n\nAll fields are empty\n\n\n"
					break
				}

				// TODO: must set the conclussion message before this
				m.Nav.SetNewMenu(config.IRQ_CONCLUSSION_VIEW_ID)

				// TODO: this needs to return a tea.Cmd (or maybe not)
				// TODO: Apply the changes call the kcmdline funcs
			}

			// Cycle indexes
			if key.Matches(msg, m.keys.Up) {
				index.Prev()
			}

			if key.Matches(msg, m.keys.Down) ||
				key.Matches(msg, m.keys.Select) {
				index.Next()
			}

			cmds := make([]tea.Cmd, len(m.Inputs))
			for i := 0; i <= len(m.Inputs)-1; i++ {
				if i == m.FocusIndex {
					// Set focused state
					cmds[i] = m.Inputs[i].Focus()
					m.Inputs[i].PromptStyle = styles.FocusedStyle
					m.Inputs[i].TextStyle = styles.FocusedStyle
					m.Inputs[i].Placeholder = placeholders_text[i]
					continue
				}
				// Remove focused state
				m.Inputs[i].Blur()
				m.Inputs[i].PromptStyle = styles.NoStyle
				m.Inputs[i].TextStyle = styles.NoStyle
				m.Inputs[i].Placeholder = ""
			}

		default:
			log.Printf("key pressed: %v", msg.String())
		}
	}
	for i := range m.Inputs {
		m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

// log.Println("---(kcmdlineMenuUpdate - start")
// cmds := make([]tea.Cmd, len(m.Inputs))

// totalIrqItems := len(m.Inputs) + 2
// index := cmp.NewNavigation(&m.FocusIndex, &totalIrqItems)
