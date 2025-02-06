package ui

import (
	"log"

	"github.com/canonical/rt-conf/src/kcmd"
	cmp "github.com/canonical/rt-conf/src/ui/components"
	"github.com/canonical/rt-conf/src/ui/config"
	"github.com/canonical/rt-conf/src/ui/styles"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	isolcpusIndex = iota
	nohzIndex
	nohzFullIndex
	kthreadsCPUsIndex
	irqaffinityIndex
	applyButtonIndex
	backButtonIndex
)

// TODO: fix the problem with the j,k keys being logged
func (m *Model) kcmdlineMenuUpdate(msg tea.KeyMsg) tea.Cmd {
	log.Println("---(kcmdlineMenuUpdate - start")
	cmds := make([]tea.Cmd, len(m.kcmd.Inputs))

	totalIrqItems := len(m.kcmd.Inputs) + 2
	index := cmp.NewNavigation(&m.kcmd.FocusIndex, &totalIrqItems)

	switch {

	case key.Matches(msg, m.kcmd.keys.Select),
		key.Matches(msg, m.kcmd.keys.Up),
		key.Matches(msg, m.kcmd.keys.Down),
		key.Matches(msg, m.kcmd.keys.Left),
		key.Matches(msg, m.kcmd.keys.Right):

		// Handle navigation between the buttons
		if m.kcmd.FocusIndex == applyButtonIndex &&
			key.Matches(msg, m.kcmd.keys.Left) {
			index.Prev()
		}
		if m.kcmd.FocusIndex == applyButtonIndex &&
			key.Matches(msg, m.kcmd.keys.Right) {
			index.Next()
		}
		if m.kcmd.FocusIndex == backButtonIndex &&
			key.Matches(msg, m.kcmd.keys.Right) {
			index.Next()
		}
		if m.kcmd.FocusIndex == backButtonIndex &&
			key.Matches(msg, m.kcmd.keys.Left) {
			index.Prev()
		}

		// log.Println("focusIndex on Update: ", m.kcmdFocusIndex)
		// Validate the inputs

		valid := m.AreValidInputs()
		log.Println("isValid: ", valid)

		// Handle [ Back ] button
		if m.kcmd.FocusIndex == backButtonIndex &&
			key.Matches(msg, m.kcmd.keys.Select) {
			// log.Println("pressed [ BACK ]: Back to main menu")
			m.nav.PrevMenu()
		}

		// Did the user press enter while the apply button was focused?
		// TODO: improve mapping of len(m.inputs) to the apply button
		if key.Matches(msg, m.kcmd.keys.Select) &&
			m.kcmd.FocusIndex == len(m.kcmd.Inputs) {

			// log.Println("Apply changes")

			valid := m.AreValidInputs()

			if !valid {
				break
			}

			m.iConf.Data.KernelCmdline.IsolCPUs = m.kcmd.Inputs[isolcpusIndex].Value()

			m.iConf.Data.KernelCmdline.Nohz = m.kcmd.Inputs[nohzIndex].Value()

			m.iConf.Data.KernelCmdline.NohzFull = m.kcmd.Inputs[nohzFullIndex].Value()

			m.iConf.Data.KernelCmdline.KthreadCPUs = m.kcmd.Inputs[kthreadsCPUsIndex].Value()

			m.iConf.Data.KernelCmdline.IRQaffinity = m.kcmd.Inputs[irqaffinityIndex].Value()

			msgs, err := kcmd.ProcessKcmdArgs(&m.iConf)
			if err != nil {
				m.errorMsg = "Failed to process kernel cmdline args: " +
					err.Error()
				break
			}

			m.logMsg = msgs
			m.renderLog = true
			m.nav.SetNewMenu(config.KCMD_CONCLUSSION_VIEW_ID)

			// TODO: this needs to return a tea.Cmd (or maybe not)
			// TODO: Apply the changes call the kcmdline funcs
		}

		// Cycle indexes
		if key.Matches(msg, m.kcmd.keys.Up) {
			// m.PrevIndex(&m.kcmdFocusIndex, m.kcmdInputs)
			index.Prev()
		}

		if key.Matches(msg, m.kcmd.keys.Down) ||
			key.Matches(msg, m.kcmd.keys.Select) {
			// m.NextIndex(&m.kcmdFocusIndex, m.kcmdInputs)
			index.Next()
		}

		cmds := make([]tea.Cmd, len(m.kcmd.Inputs))
		for i := 0; i <= len(m.kcmd.Inputs)-1; i++ {
			if i == m.kcmd.FocusIndex {
				// Set focused state
				cmds[i] = m.kcmd.Inputs[i].Focus()
				m.kcmd.Inputs[i].PromptStyle = styles.FocusedStyle
				m.kcmd.Inputs[i].TextStyle = styles.FocusedStyle
				m.kcmd.Inputs[i].Placeholder = placeholders_text[i]
				continue
			}
			// Remove focused state
			m.kcmd.Inputs[i].Blur()
			m.kcmd.Inputs[i].PromptStyle = styles.NoStyle
			m.kcmd.Inputs[i].TextStyle = styles.NoStyle
			m.kcmd.Inputs[i].Placeholder = ""
		}
	}
	for i := range m.kcmd.Inputs {
		m.kcmd.Inputs[i], cmds[i] = m.kcmd.Inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m *Model) kcmdlineConclussionUpdate(msg tea.KeyMsg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.kcmd.Inputs))
	log.Println("(kcmdlineConclussionUpdate - start")

	switch {
	/* If the user press enter on the log view,
	go back to the previous menu */
	case key.Matches(msg, m.kcmd.keys.Select):
		m.renderLog = false
		m.nav.PrevMenu()

	default:
		log.Println("(default) Invalid key pressed: ", msg.String())
		log.Println("Menu navigation stack: ", m.nav.PrintMenuStack())

	}
	return tea.Batch(cmds...)
}
