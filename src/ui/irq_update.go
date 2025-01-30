package ui

import (
	"log"

	"github.com/canonical/rt-conf/src/kcmd"
	"github.com/canonical/rt-conf/src/ui/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// ** NOTE: The index navigation needs to be dynamic now
// const (
// 	irqCPUsIndex = iota
// 	irqNumberIndex
// 	irqActionIndex
// 	irqChipNameIndex
// 	irqNameIndex
// 	irqTypeIndex
// )

// TODO: Functions to generate new entries for IRQ affinity menu
func NewIRQTextInputs() textinput.Model {
	panic("not implemented")
}

func (m *Model) updateIRQMenu(msg tea.KeyMsg) tea.Cmd {

	switch {
	case key.Matches(msg, m.keys.Help):
		// Once detected the key "?" toggle the help message
		// but first disable the text input
		var cmd tea.Cmd

		// Checking for overflow since the buttons aren't text inputs
		if m.focusIndex < applyButtonIndex {
			m.kcmdInputs[m.focusIndex].Blur()
		}
		m.help.ShowAll = !m.help.ShowAll

		if m.focusIndex < applyButtonIndex {
			cmd = m.kcmdInputs[m.focusIndex].Focus()
		}
		return cmd

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

		log.Println("focusIndex on Update: ", m.focusIndex)

		// Validate the inputs
		log.Println("isValid: ", m.AreValidInputs())

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
			m.focusIndex == len(m.irqInputs) {

			log.Println("Apply changes")

			valid := m.AreValidInputs()

			if !valid {
				break
			}

			m.iConf.Data.KernelCmdline.IsolCPUs = m.irqInputs[isolcpusIndex].Value()

			m.iConf.Data.KernelCmdline.Nohz = m.irqInputs[nohzIndex].Value()

			m.iConf.Data.KernelCmdline.NohzFull = m.irqInputs[nohzFullIndex].Value()

			m.iConf.Data.KernelCmdline.KthreadCPUs = m.irqInputs[kthreadsCPUsIndex].Value()

			m.iConf.Data.KernelCmdline.IRQaffinity = m.irqInputs[irqaffinityIndex].Value()

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

		// Cycle indexes
		if key.Matches(msg, m.keys.Up) {
			m.PrevIndex()
		}

		if key.Matches(msg, m.keys.Down) ||
			key.Matches(msg, m.keys.Select) {
			m.NextIndex()
		}

		cmds := make([]tea.Cmd, len(m.kcmdInputs))
		for i := 0; i <= len(m.kcmdInputs)-1; i++ {
			if i == m.focusIndex {
				// Set focused state
				cmds[i] = m.kcmdInputs[i].Focus()
				m.kcmdInputs[i].PromptStyle = styles.FocusedStyle
				m.kcmdInputs[i].TextStyle = styles.FocusedStyle
				m.kcmdInputs[i].Placeholder = placeholders_text[i]
				continue
			}
			// Remove focused state
			m.kcmdInputs[i].Blur()
			m.kcmdInputs[i].PromptStyle = styles.NoStyle
			m.kcmdInputs[i].TextStyle = styles.NoStyle
			m.kcmdInputs[i].Placeholder = ""
		}

		return tea.Batch(cmds...)
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

	return nil
}
