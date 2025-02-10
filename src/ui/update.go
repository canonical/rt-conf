package ui

import (
	"log"
	"strings"

	"github.com/canonical/rt-conf/src/data"
	"github.com/canonical/rt-conf/src/interrupts"
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

var placeholders_irq = []string{
	config.IrqFilterPlaceholder,
	config.CpuListPlaceholder,
}

const (
	irqFilterIndex = iota
	cpuListIndex
	addBtnIndex
	cancelBtnIndex
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Println("\n------------- UPDATE -----------------")
	// log.Println("(MAIN UPDATE) Current menu: ", config.Menu[m.Nav.GetCurrMenu()])
	// log.Println("MENU STACK: ", m.Nav.PrintMenuStack())
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

		// Main Menu
		m.main.list.SetSize(msg.Width-h, msg.Height-v)

		// Kcmdline
		m.kcmd.help.Width = msg.Width
		m.kcmd.Width = msg.Width - h
		m.kcmd.Height = msg.Height - v
		m.kcmd.concl.Height = msg.Height - v
		m.kcmd.concl.Width = msg.Width - h

		// IRQ
		m.irq.list.SetSize(msg.Width-h, msg.Height-v)
		m.irq.Width = msg.Width - h
		m.irq.Height = msg.Height - v
		m.irq.irq.width = msg.Width - h
		m.irq.irq.height = msg.Height - v
		m.irq.concl.Width = msg.Width - h
		m.irq.concl.Height = msg.Height - v

	case tea.KeyMsg:
		// log.Printf("Key pressed: %v", msg.String())
		switch {
		case key.Matches(msg, m.main.keys.Back):
			m.Nav.PrevMenu()
			return m, nil
		// This is genertic for all menus

		case key.Matches(msg, m.main.keys.Quit):
			if msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
			return m, nil

		case key.Matches(msg, m.main.keys.Home):
			m.Nav.ReturnToMainMenu()

		}
	}

	activeMenu := m.GetActiveMenu()
	var cmd tea.Cmd
	// Calling the Update() function of the active menu
	m.currMenu, cmd = activeMenu.Update(msg)
	cmds = append(cmds, cmd)

	// TODO: improve this workaround
	// This is an excecption for the IRQ Add/Editview so the 'q' key can be used
	if m.Nav.GetCurrMenu() == config.IRQ_ADD_EDIT_VIEW_ID {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q":
				return m, tea.Batch(cmds...)
			}
		}
	}

	// Update main menu list view
	newListModel, cmd := m.main.list.Update(msg)
	m.main.list = newListModel
	cmds = append(cmds, cmd)

	// Update Kcmdline list view // TODO: implement this
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

func (m *IRQMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("---IRQMenuUpdate: ")

	switch msg := msg.(type) {
	case IRQRuleMsg:
		log.Println("<><><><>(RECEIVED) IRQ Affinity Rule: <><><><>")
		irqRule := msg.IRQAffinityRule
		// items := m.list.Items()
		newFilter := config.PrefixIRQFilter + irqRule.filter
		newCpuList := config.PrefixCpuList + irqRule.cpulist
		// log.Println("->>NewFilter:  ", newFilter)
		// log.Println("->>NewCpuList:  ", newCpuList)

		var cmd tea.Cmd
		item := IRQAffinityRule{
			rule:    irqRule.rule,
			filter:  newFilter,
			cpulist: newCpuList,
		}
		if msg.edited {
			// log.Println("->> Edit mode")
			// log.Println("->> Index: ", irqRule.index)
			cmd = m.list.SetItem(irqRule.index, item)
			m.rules[irqRule.index] = irqRule.rule
		} else {
			index := len(m.list.Items())
			cmd = m.list.InsertItem(index, item)
			m.rules = append(m.rules, irqRule.rule)
		}
		return m, cmd

	case tea.KeyMsg:
		// log.Printf("key pressed: %v", msg.String())

		switch {

		// case key.Matches(msg, m.keys.Help):
		// 	log.Println("->>> Help key pressed <<<-")

		case key.Matches(msg, m.keys.Apply):
			// log.Println("-> Entry: Apply changes")
			num := len(m.rules)
			if len(m.rules) == len(m.list.Items()) {
				log.Println(">>> The size is coeherent")
			} else {
				log.Println("[ERROR] >>> The size is NOT coeherent")
				// panic("The size is NOT coeherent")
				// os.Exit(1)
			}
			if num > 0 {
				m.concl.num = num
				m.concl.logMsg =
					"IRQ affinity rules are aplied to the system"

				var cfg data.InternalConfig
				cfg.Data.Interrupts = m.rules

				err := interrupts.ApplyIRQConfig(&cfg)
				if err != nil {
					log.Println("ERROR: ", err)
					m.concl.num = 0
					m.concl.logMsg = m.concl.logMsg + "\nERROR: " + err.Error()
				}
				m.Nav.SetNewMenu(config.IRQ_CONCLUSION_VIEW_ID)
			}

		case key.Matches(msg, m.keys.Add):
			// log.Println("-> Entry: Add new IRQ rule")
			// TODO: implement cleanup of AddEditMenu fields
			cmd := func() tea.Msg {
				return StartNewIRQAffinityRule{editMode: false}
			}
			m.Nav.SetNewMenu(config.IRQ_ADD_EDIT_VIEW_ID)
			return m, cmd

		case key.Matches(msg, m.keys.Remove):
			// log.Println("-> Entry: Remove IRQ rule")
			idx := m.list.Index()
			m.rules = append(m.rules[:idx], m.rules[idx+1:]...)
			m.list.RemoveItem(idx)

		case key.Matches(msg, m.irq.keys.Select):
			// log.Printf("confirmed select pressed")
			// log.Printf("previous menu: %v", config.Menu[m.Nav.GetCurrMenu()])
			cpulist := strings.TrimPrefix(
				m.list.SelectedItem().(IRQAffinityRule).cpulist,
				config.PrefixCpuList,
			)
			filter := strings.TrimPrefix(
				m.list.SelectedItem().(IRQAffinityRule).filter,
				config.PrefixIRQFilter,
			)
			cmd := func() tea.Msg {
				return StartNewIRQAffinityRule{
					editMode: true,
					index:    m.list.Index(),
					cpulist:  cpulist,
					filter:   filter,
				}
			}
			m.Nav.SetNewMenu(config.IRQ_ADD_EDIT_VIEW_ID)

			// log.Println("Current menu: ", config.Menu[m.Nav.GetCurrMenu()])
			return m, cmd

		}
	}
	return m, nil
}

func (m *IRQConclusion) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	cmds := make([]tea.Cmd, len(m.Inputs), len(m.Inputs)+4)

	totalIrqItems := len(m.Inputs) + 2
	index := cmp.NewNavigation(&m.FocusIndex, &totalIrqItems)
	// log.Println("\tFocus index: ", m.FocusIndex)
	// log.Println("\tAddres: ", &m.FocusIndex)
	// log.Printf("---IRQAddEditMenuUpdate: ")
	// log.Printf("Current menu: %v", config.Menu[m.Nav.GetCurrMenu()])

	// TODO: here it should be implemented the logic to add or edit IRQs
	switch msg := msg.(type) {

	case StartNewIRQAffinityRule:
		// log.Println("<><><> StartNewIRQAffinityRule")
		// log.Println(">>> Edit mode: ", msg.editMode)
		// log.Println(">>> Index: ", msg.index)

		m.editMode = msg.editMode
		m.editIndex = msg.index
		m.FocusIndex = 0
		m.Inputs[irqFilterIndex].Focus()
		m.Inputs[irqFilterIndex].PromptStyle = styles.FocusedStyle
		m.Inputs[irqFilterIndex].PromptStyle = styles.FocusedStyle
		m.Inputs[irqFilterIndex].TextStyle = styles.FocusedStyle
		m.Inputs[irqFilterIndex].Placeholder = placeholders_irq[irqFilterIndex]
		if !msg.editMode {
			m.Inputs[irqFilterIndex].SetValue("")
			m.Inputs[cpuListIndex].SetValue("")
		} else {
			m.Inputs[irqFilterIndex].SetValue(msg.filter)
			m.Inputs[cpuListIndex].SetValue(msg.cpulist)
		}

	case tea.KeyMsg:
		var isValid bool
		switch {
		// TODO: find a way to allow use of 'q' key on the text inputs
		case key.Matches(msg, m.keys.Quit):
			if msg.String() == "ctrl+c" {
				return m, tea.Quit
			}

		case key.Matches(msg, m.keys.Help):
			var cmd tea.Cmd
			if m.FocusIndex < addBtnIndex {
				m.Inputs[m.FocusIndex].Blur()
			}
			m.help.ShowAll = !m.help.ShowAll
			if m.FocusIndex < addBtnIndex {
				cmd = m.Inputs[m.FocusIndex].Focus()
			}
			return m, tea.Sequence(cmd)

		case key.Matches(msg, m.keys.Up):
			m.RunInputValidation()
			index.Prev()

		case key.Matches(msg, m.keys.Down):
			m.RunInputValidation()
			index.Next()

		case key.Matches(msg, m.keys.Left):
			if m.FocusIndex == addBtnIndex {
				index.Prev()
			} else if m.FocusIndex == cancelBtnIndex {
				index.Prev()
			}

		case key.Matches(msg, m.keys.Right):
			if m.FocusIndex == addBtnIndex {
				index.Next()
			} else if m.FocusIndex == cancelBtnIndex {
				index.Next()
			}

		case key.Matches(msg, m.keys.Select):
			isValid = m.AreValidInputs()
			// log.Println(">>> Select key pressed")

			// TODO: fix bug where it's necessary to press twice the buttons
			// Handle [ Cancel ] button
			if m.FocusIndex == cancelBtnIndex {
				// log.Println(">>> Cancel changes")
				m.Nav.PrevMenu()
				break
			}

			if m.FocusIndex < len(m.Inputs) /*addBtnIndex*/ {
				index.Next()
				break
			}

			// Did the user press enter while the [ ADD ] button was focused?
			if m.FocusIndex == addBtnIndex {
				log.Println(">>> Apply changes | isValid: ", isValid)

				var empty int
				for i := range m.Inputs {
					v := m.Inputs[i].Value()
					if v == "" {
						empty++
					}
				}
				if empty == len(m.Inputs) {
					isValid = false
					m.errorMsg = "\nAll fields are empty\n"
					break
				}
				if empty > 0 {
					isValid = false
					m.errorMsg = "\nA field is empty\n"
					break
				}

				if !isValid {
					break
				}

				rawFilter := m.Inputs[irqFilterIndex].Value()
				cpulist := m.Inputs[cpuListIndex].Value()

				irqFilter, _ := ParseIRQFilter(rawFilter)
				// if err != nil {
				// 	m.errorMsg = "ERROR: " + err.Error() + "\n"
				// 	break
				// }

				irqAffinityRule := IRQAffinityRule{
					rule: data.IRQTunning{
						Filter: irqFilter,
						CPUs:   cpulist,
					},
					filter:  rawFilter,
					cpulist: cpulist,
					edited:  m.editMode,
					index:   m.editIndex,
				}
				cmd := func() tea.Msg { return IRQRuleMsg{irqAffinityRule} }
				cmds = append(cmds, cmd)

				// TODO once validated, it must:
				// ** 1. Parse the values for the IRQFilter struct
				// ** 1.1 The func for parsing of IRQFilter  must be implemented
				// ** 1.2 The cpuList value can be copy directly

				// ** 2. Create a tea.Cmd message to send the IRQFilter struct
				// ** 2.1 The tea.Cmd must be added into the tea.Batch(cmds...)
				m.Nav.PrevMenu()

			}

		default:
			// log.Printf("key pressed: %v", msg.String())
		}
	}

	for i := 0; i <= len(m.Inputs)-1; i++ {
		if i == m.FocusIndex {
			// Set focused state
			cmds[i] = m.Inputs[i].Focus()
			m.Inputs[i].PromptStyle = styles.FocusedStyle
			m.Inputs[i].TextStyle = styles.FocusedStyle
			m.Inputs[i].Placeholder = placeholders_irq[i]
			continue
		}
		// Remove focused state
		m.Inputs[i].Blur()
		m.Inputs[i].PromptStyle = styles.NoStyle
		m.Inputs[i].TextStyle = styles.NoStyle
		m.Inputs[i].Placeholder = ""
	}

	for i := range m.Inputs {
		m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}
