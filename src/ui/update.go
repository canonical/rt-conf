package ui

import (
	"log"
	"strconv"

	"github.com/canonical/rt-conf/src/kcmd"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// TODO: FOCUS ON THE [ APPLY ] FUNCTIONALITY
// TODO:

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmdInput tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.width = msg.Width - h
		m.height = msg.Height - v

		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.help.Width = msg.Width

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.

		if m.currMenu != mainMenu {

			// TODO: Move this logic to a func specific for kcmdlineMenu
			switch msg.String() {
			case "ctrl+c", "esc":
				return m, tea.Quit

			case "g":
				m.currMenu = mainMenu
				return m, nil

			// Change cursor mode
			case "ctrl+r":
				m.cursorMode++
				if m.cursorMode > cursor.CursorHide {
					m.cursorMode = cursor.CursorBlink
				}
				cmds := make([]tea.Cmd, len(m.inputs))
				for i := range m.inputs {
					cmds[i] = m.inputs[i].Cursor.SetMode(m.cursorMode)
				}
				return m, tea.Batch(cmds...)

			// Set focus to next input
			case "tab", "shift+tab", "enter", "up", "down":
				s := msg.String()

				log.Println("focusIndex on Update: ", m.focusIndex)

				isValid := m.Validation()
				log.Println("isValid: ", isValid)

				// Did the user press enter while the submit button was focused?
				// If so, exit.
				if s == "enter" && m.focusIndex == len(m.inputs) {
					log.Println("Apply changes")

					valid := m.Validation()

					dyntickIdle, err := strconv.ParseBool(m.inputs[enableDynticks].Value())
					if err != nil {
						log.Printf("Failed to parse dyntick idle value: %v", err)
						break
					}
					if !valid {
						break
					}
					// TODO: Improve this logic
					m.iconf.Data.KernelCmdline.DyntickIdle = dyntickIdle
					m.iconf.Data.KernelCmdline.IsolCPUs = m.inputs[isolatecpus].Value()
					m.iconf.Data.KernelCmdline.AdaptiveTicks = m.inputs[enableDynticks].Value()

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
				if s == "up" || s == "shift+tab" {
					m.PrevIndex()
				} else {
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
						m.inputs[i].PromptStyle = focusedStyle
						m.inputs[i].TextStyle = focusedStyle
						continue
					}
					// Remove focused state
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = noStyle
					m.inputs[i].TextStyle = noStyle
				}

				return m, tea.Batch(cmds...)
			}
		}
	}

	cmdMainMenu := m.updateMainMenu(msg)
	cmdInput = m.updateInputs(msg)
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)
	cmds = append(cmds, cmdInput)
	cmds = append(cmds, cmdMainMenu)

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
			case key.Matches(msg, m.keys.toggleSpinner):
				cmd := m.list.ToggleSpinner()
				cmds = append(cmds, cmd)

			case key.Matches(msg, m.keys.goHome):
				m.currMenu = mainMenu

			case key.Matches(msg, m.keys.selectMenu):
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

			case key.Matches(msg, m.keys.toggleHelpMenu):
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
