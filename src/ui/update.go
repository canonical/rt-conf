package ui

import (
	"log"

	"github.com/canonical/rt-conf/src/kcmd"
	"github.com/canonical/rt-conf/src/ui/styles"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// TODO: FOCUS ON THE [ APPLY ] FUNCTIONALITY

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
				m.inputs[m.focusIndex].Blur()
				m.help.ShowAll = !m.help.ShowAll
				return m, m.inputs[m.focusIndex].Focus()

			case key.Matches(msg, m.keys.CursorMode):
				m.cursorMode++
				if m.cursorMode > cursor.CursorHide {
					m.cursorMode = cursor.CursorBlink
				}
				cmds := make([]tea.Cmd, len(m.inputs))
				for i := range m.inputs {
					cmds[i] = m.inputs[i].Cursor.SetMode(m.cursorMode)
				}
				return m, tea.Batch(cmds...)

			case key.Matches(msg, m.keys.Select),
				key.Matches(msg, m.keys.Up),
				key.Matches(msg, m.keys.Down):
				// s := msg.String()
				log.Println("focusIndex on Update: ", m.focusIndex)

				isValid := m.Validation()
				log.Println("isValid: ", isValid)

				// Did the user press enter while the submit button was focused?
				// If so, exit.
				if key.Matches(msg, m.keys.Select) &&
					m.focusIndex == len(m.inputs) {

					log.Println("Apply changes")

					valid := m.Validation()

					var dyntickIdle bool

					// TODO: move this away from here
					if m.inputs[enableDynticks].Value() == "y" || m.inputs[enableDynticks].Value() == "Y" {
						dyntickIdle = true
					} else if m.inputs[enableDynticks].Value() == "n" || m.inputs[enableDynticks].Value() == "N" {
						dyntickIdle = false
					} else {
						m.errorMsg = "ERROR: expected Yes or No value (y|n) got: " + m.inputs[enableDynticks].Value()
						break
					}

					if !valid {
						break
					}
					// TODO: Improve this logic
					m.iconf.Data.KernelCmdline.DyntickIdle = dyntickIdle
					m.iconf.Data.KernelCmdline.IsolCPUs = m.inputs[isolatecpus].Value()
					m.iconf.Data.KernelCmdline.AdaptiveTicks = m.inputs[adaptiveCPUs].Value()

					msgs, err := kcmd.ProcessKcmdArgs(&m.iconf)
					if err != nil {
						m.errorMsg = "Failed to process kernel cmdline args: " + err.Error()
						break
					}
					m.infoMsg = "\n" // Doesn't show the info message

					m.logMsg = msgs
					m.renderLog = true

					// TODO: this needs to return a tea.Cmd (or maybe not)

					// TODO: Apply the changes call the kcmdline funcs
					// TODO: Present a new view (maybe a new block with text)
				}

				log.Println("Focus index is: ", m.focusIndex)
				if m.focusIndex < len(m.inputs) {
					log.Println("Value: ", m.inputs[m.focusIndex].Value())
				}

				// Cycle indexes
				if key.Matches(msg, m.keys.Up) {
					m.PrevIndex()
				}

				if key.Matches(msg, m.keys.Down) || key.Matches(msg, m.keys.Select) {
					m.NextIndex()
				}

				// if m.focusIndex > len(m.inputs) {
				// 	m.focusIndex = 0
				// } else if m.focusIndex < 0 {
				// 	m.focusIndex = len(m.inputs)
				// }

				cmds := make([]tea.Cmd, len(m.inputs))
				for i := 0; i <= len(m.inputs)-1; i++ {
					if i == m.focusIndex {
						// Set focused state
						cmds[i] = m.inputs[i].Focus()
						m.inputs[i].PromptStyle = styles.FocusedStyle
						m.inputs[i].TextStyle = styles.FocusedStyle
						continue
					}
					// Remove focused state
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = styles.NoStyle
					m.inputs[i].TextStyle = styles.NoStyle
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
		// Only text inputs with Focus() set will respond, so it's safe to simply
		// update all of them here without any further logic.
		for i := range m.inputs {
			m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
		}
	}

	return tea.Batch(cmds...)
}
