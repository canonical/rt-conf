package ui

import (
	"log"

	"github.com/canonical/rt-conf/src/data"
	cmp "github.com/canonical/rt-conf/src/ui/components"
	"github.com/canonical/rt-conf/src/ui/styles"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

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
