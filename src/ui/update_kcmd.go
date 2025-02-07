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
func (m KcmdlineMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	log.Println("---(kcmdlineMenuUpdate - start")
	cmds := make([]tea.Cmd, len(m.Inputs))

	totalIrqItems := len(m.Inputs) + 2
	index := cmp.NewNavigation(&m.FocusIndex, &totalIrqItems)

	switch msg := msg.(type) {
	// case tea.WindowSizeMsg:
	// 	m.conclussion.Width = msg.Width
	// 	m.conclussion.Height = msg.Height

	// 	m.help.Width = msg.Width

	case tea.KeyMsg:
		switch {

		case key.Matches(msg, m.keys.Select),
			key.Matches(msg, m.keys.Up),
			key.Matches(msg, m.keys.Down),
			key.Matches(msg, m.keys.Left),
			key.Matches(msg, m.keys.Right):

			// Handle navigation between the buttons
			if m.FocusIndex == applyButtonIndex &&
				key.Matches(msg, m.keys.Left) {
				index.Prev()
			}
			if m.FocusIndex == applyButtonIndex &&
				key.Matches(msg, m.keys.Right) {
				index.Next()
			}
			if m.FocusIndex == backButtonIndex &&
				key.Matches(msg, m.keys.Right) {
				index.Next()
			}
			if m.FocusIndex == backButtonIndex &&
				key.Matches(msg, m.keys.Left) {
				index.Prev()
			}

			// log.Println("focusIndex on Update: ", m.kcmdFocusIndex)
			// Validate the inputs

			valid := m.AreValidInputs()
			// log.Println("isValid: ", valid)

			// Handle [ Back ] button
			if m.FocusIndex == backButtonIndex &&
				key.Matches(msg, m.keys.Select) {
				// log.Println("pressed [ BACK ]: Back to main menu")
				m.Nav.PrevMenu()
			}

			// Did the user press enter while the apply button was focused?
			// TODO: improve mapping of len(m.inputs) to the apply button
			if key.Matches(msg, m.keys.Select) &&
				m.FocusIndex == len(m.Inputs) && valid {

				// log.Println("Apply changes")

				valid := m.AreValidInputs()

				if !valid {
					break
				}
				var empty int
				for i := range m.Inputs {
					v := m.Inputs[i].Value()
					if v == "" {
						empty++
					}
				}
				if empty == len(m.Inputs) {
					m.errorMsg = "\n\nAll fields are empty\n\n\n"
					break
				}

				m.iConf.Data.KernelCmdline.IsolCPUs = m.Inputs[isolcpusIndex].Value()

				m.iConf.Data.KernelCmdline.Nohz = m.Inputs[nohzIndex].Value()

				m.iConf.Data.KernelCmdline.NohzFull = m.Inputs[nohzFullIndex].Value()

				m.iConf.Data.KernelCmdline.KthreadCPUs = m.Inputs[kthreadsCPUsIndex].Value()

				m.iConf.Data.KernelCmdline.IRQaffinity = m.Inputs[irqaffinityIndex].Value()

				msgs, err := kcmd.ProcessKcmdArgs(&m.iConf)
				if err != nil {
					m.errorMsg = "Failed to process kernel cmdline args: " +
						err.Error()
					break
				}

				m.conclussion.logMsg = msgs
				m.conclussion.renderLog = true
				m.Nav.SetNewMenu(config.KCMD_CONCLUSSION_VIEW_ID)

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
		}
	}
	for i := range m.Inputs {
		m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m *KcmdlineConclussion) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	log.Println("(kcmdlineConclussionUpdate - start")

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:
		switch {
		/* If the user press enter on the log view,
		go back to the previous menu */
		case key.Matches(msg, m.keys.Select):
			// m.renderLog = false
			m.Nav.PrevMenu()
		default:
			log.Println("(default) Invalid key pressed: ", msg.String())
			log.Println("Menu navigation stack: ", m.Nav.PrintMenuStack())
		}
	}
	return m, cmd
}
