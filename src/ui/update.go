package ui

import (
	"log"

	"github.com/canonical/rt-conf/src/kcmd"
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
	applyButtonIndex
	backButtonIndex
)

const cpuListPlaceholder = "Enter a CPU list like: 4-n or 3-5 or 2,4,5 "

var placeholders_text = []string{
	cpuListPlaceholder,
	"Enter y or n",
	cpuListPlaceholder,
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmdInput tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// TODO: check if this will not break the app
		h, v := styles.AppStyle.GetFrameSize()

		m.width = msg.Width - h
		m.height = msg.Height - v

		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.help.Width = msg.Width

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.

		if m.currMenu != mainMenu {
			switch {
			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit

			case key.Matches(msg, m.keys.goHome):
				m.currMenu = mainMenu

			case key.Matches(msg, m.keys.Help):
				// Once detected the key "?" toggle the help message
				// but first disable the text input
				var cmd tea.Cmd

				/* Checking for overflow since the apply button isn't a
				text input */
				if m.focusIndex < applyButtonIndex {
					m.inputs[m.focusIndex].Blur()
				}
				m.help.ShowAll = !m.help.ShowAll

				if m.focusIndex < applyButtonIndex {
					cmd = m.inputs[m.focusIndex].Focus()
				}
				return m, cmd

			case key.Matches(msg, m.keys.Select),
				key.Matches(msg, m.keys.Up),
				key.Matches(msg, m.keys.Down),
				key.Matches(msg, m.keys.Left),
				key.Matches(msg, m.keys.Right):

				// Handle navigation between the buttons
				if m.focusIndex == applyButtonIndex &&
					key.Matches(msg, m.keys.Left) {
					m.focusIndex = backButtonIndex
				}
				if m.focusIndex == backButtonIndex &&
					key.Matches(msg, m.keys.Up) {
					m.focusIndex--
				}
				if m.focusIndex == backButtonIndex &&
					key.Matches(msg, m.keys.Right) {
					m.focusIndex = applyButtonIndex
				}

				// s := msg.String()
				log.Println("focusIndex on Update: ", m.focusIndex)

				// Validate the inputs
				isValid := m.Validation()
				log.Println("isValid: ", isValid)

				// Handle [ Back ] button
				if m.focusIndex == backButtonIndex &&
					key.Matches(msg, m.keys.Select) {
					m.currMenu = mainMenu
				}

				/* If the user press enter on the log view,
				go back to the previous menu */
				if m.renderLog && key.Matches(msg, m.keys.Select) {
					m.renderLog = false
					m.currMenu = kcmdlineMenu
				}

				// Did the user press enter while the apply button was focused?
				// TODO: improve mapping of len(m.inputs) to the apply button
				if key.Matches(msg, m.keys.Select) &&
					m.focusIndex == len(m.inputs) {

					log.Println("Apply changes")

					valid := m.AreValidInputs()

					if !valid {
						break
					}

					// TODO: Improve this logic
					nohzVal := m.inputs[nohzIndex].Value()
					if nohzVal == "y" || nohzVal == "Y" {
						m.iConf.Data.KernelCmdline.DyntickIdle = true
					} else {
						m.iConf.Data.KernelCmdline.DyntickIdle = false
					}

					m.iConf.Data.KernelCmdline.IsolCPUs =
						m.inputs[isolcpusIndex].Value()

					m.iConf.Data.KernelCmdline.AdaptiveTicks =
						m.inputs[nohzFullIndex].Value()

					msgs, err := kcmd.ProcessKcmdArgs(&m.iConf)
					if err != nil {
						m.errorMsg = "Failed to process kernel cmdline args: " +
							err.Error()
						break
					}

					m.logMsg = msgs
					m.renderLog = true

					// TODO: this needs to return a tea.Cmd (or maybe not)

					// TODO: Apply the changes call the kcmdline funcs
				}

				log.Println("Focus index is: ", m.focusIndex)
				if m.focusIndex < len(m.inputs) {
					log.Println("Value: ", m.inputs[m.focusIndex].Value())
				}

				// Cycle indexes
				if key.Matches(msg, m.keys.Up) {
					m.PrevIndex()
				}

				if key.Matches(msg, m.keys.Down) ||
					key.Matches(msg, m.keys.Select) {
					m.NextIndex()
				}

				cmds := make([]tea.Cmd, len(m.inputs))
				for i := 0; i <= len(m.inputs)-1; i++ {
					if i == m.focusIndex {
						// Set focused state
						cmds[i] = m.inputs[i].Focus()
						m.inputs[i].PromptStyle = styles.FocusedStyle
						m.inputs[i].TextStyle = styles.FocusedStyle
						m.inputs[i].Placeholder = placeholders_text[i]
						continue
					}
					// Remove focused state
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = styles.NoStyle
					m.inputs[i].TextStyle = styles.NoStyle
					m.inputs[i].Placeholder = ""
				}

				return m, tea.Batch(cmds...)

			}
		}
	}

	cmdMainMenu := m.updateMainMenu(msg)
	cmdInput = m.updateInputs(msg)
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd, cmdInput, cmdMainMenu)

	return m, tea.Batch(cmds...)
}

func (m *Model) updateMainMenu(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	if m.currMenu == mainMenu {
		switch msg := msg.(type) {
		case tea.KeyMsg:

			if m.list.FilterState() == list.Filtering {
				break
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

		}

	}

	return tea.Batch(cmds...)
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	if m.currMenu == kcmdlineMenu {
		// Only text inputs with Focus() set will respond, so it's safe
		// to simply update all of them here without any further logic.
		for i := range m.inputs {
			m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
		}
	}

	return tea.Batch(cmds...)
}
