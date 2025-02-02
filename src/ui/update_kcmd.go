package ui

import (
	"log"

	"github.com/canonical/rt-conf/src/kcmd"
	"github.com/canonical/rt-conf/src/ui/styles"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// TODO: fix the problem with the j,k keys being logged
func (m *Model) updateKcmdlineMenu(msg tea.KeyMsg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.kcmdInputs))

	switch {

	case key.Matches(msg, m.keys.Select),
		key.Matches(msg, m.keys.Up),
		key.Matches(msg, m.keys.Down),
		key.Matches(msg, m.keys.Left),
		key.Matches(msg, m.keys.Right):

		// Handle navigation between the buttons
		if m.kcmdFocusIndex == applyButtonIndex &&
			key.Matches(msg, m.keys.Left) {
			m.kcmdFocusIndex = backButtonIndex
		}
		if m.kcmdFocusIndex == backButtonIndex &&
			key.Matches(msg, m.keys.Up) {
			m.kcmdFocusIndex--
		}
		if m.kcmdFocusIndex == backButtonIndex &&
			key.Matches(msg, m.keys.Right) {
			m.kcmdFocusIndex = applyButtonIndex
		}

		log.Println("focusIndex on Update: ", m.kcmdFocusIndex)

		// Validate the inputs
		log.Println("isValid: ", m.AreValidInputs())

		// Handle [ Back ] button
		if m.kcmdFocusIndex == backButtonIndex &&
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
			m.kcmdFocusIndex == len(m.kcmdInputs) {

			log.Println("Apply changes")

			valid := m.AreValidInputs()

			if !valid {
				break
			}

			m.iConf.Data.KernelCmdline.IsolCPUs = m.kcmdInputs[isolcpusIndex].Value()

			m.iConf.Data.KernelCmdline.Nohz = m.kcmdInputs[nohzIndex].Value()

			m.iConf.Data.KernelCmdline.NohzFull = m.kcmdInputs[nohzFullIndex].Value()

			m.iConf.Data.KernelCmdline.KthreadCPUs = m.kcmdInputs[kthreadsCPUsIndex].Value()

			m.iConf.Data.KernelCmdline.IRQaffinity = m.kcmdInputs[irqaffinityIndex].Value()

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
			m.PrevIndex(&m.kcmdFocusIndex, m.kcmdInputs)
		}

		if key.Matches(msg, m.keys.Down) ||
			key.Matches(msg, m.keys.Select) {
			m.NextIndex(&m.kcmdFocusIndex, m.kcmdInputs)
		}

		cmds := make([]tea.Cmd, len(m.kcmdInputs))
		for i := 0; i <= len(m.kcmdInputs)-1; i++ {
			if i == m.kcmdFocusIndex {
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
	}
	for i := range m.kcmdInputs {
		m.kcmdInputs[i], cmds[i] = m.kcmdInputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}
