package ui

import (
	"log"

	"github.com/canonical/rt-conf/src/ui/styles"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// ** NOTE: this is being initializaed because the IRQinputs are always in pairs
// ** so the plusBtnIndex in the beggining will be: len(m.irqInputs) which is 2
// ** But this is necessary to avoid that the plus button is focused when the
// ** menu is rendered for the first time
var plusBtnIndex = 2 // This is the initial index of the + button
var minusBtnIndex = plusBtnIndex + 1
var backBtnIndex = minusBtnIndex + 1
var applyBtnIndex = backBtnIndex + 1

// TODO: Functions to generate new entries for IRQ affinity menu
func (m *Model) NewIRQTextInputs() {
	m.irqInputs = append(m.irqInputs, newIRQtextInputs()...)
	plusBtnIndex = len(m.irqInputs)
	minusBtnIndex = plusBtnIndex + 1
	m.irqFocusIndex = plusBtnIndex
}

func (m *Model) DeleteLastIRQInput() {
	if len(m.irqInputs) > 0 {
		m.irqInputs = m.irqInputs[:len(m.irqInputs)-2]
	}
	plusBtnIndex = len(m.irqInputs)
	minusBtnIndex = plusBtnIndex + 1
	m.irqFocusIndex = minusBtnIndex
}

// TODO: add exception to keys h,j,k,l not being handled as navigation keys
func (m *Model) updateIRQMenu(msg tea.KeyMsg) tea.Cmd {
	totalIrqItems := len(m.irqInputs) + 4
	cmds := make([]tea.Cmd, len(m.irqInputs))

	log.Println("Size of IRQ inputs: ", len(m.irqInputs))
	log.Println("focusIndex on Update: ", m.irqFocusIndex)

	plusBtnIndex = len(m.irqInputs)
	minusBtnIndex = plusBtnIndex + 1
	backBtnIndex = minusBtnIndex + 1
	applyBtnIndex = backBtnIndex + 1

	dbgHelper := make(map[int]string, len(m.irqInputs)+4)
	dbgHelper[plusBtnIndex] = "Plus btn"
	dbgHelper[minusBtnIndex] = "Minus btn"
	dbgHelper[backBtnIndex] = "Back btn"
	dbgHelper[applyBtnIndex] = "Apply btn"

	log.Printf("FocusIndex: %d = %s ", m.irqFocusIndex, dbgHelper[m.irqFocusIndex])

	switch {

	case key.Matches(msg, m.keys.Select),
		key.Matches(msg, m.keys.Up),
		key.Matches(msg, m.keys.Down),
		key.Matches(msg, m.keys.Left),
		key.Matches(msg, m.keys.Right):

		// TODO: add logic for + btn
		// TODO: check weird behavior of this + button
		if m.irqFocusIndex == plusBtnIndex &&
			key.Matches(msg, m.keys.Select) {
			m.NewIRQTextInputs()
			return tea.Batch(cmds...)
		}
		if m.irqFocusIndex == minusBtnIndex &&
			key.Matches(msg, m.keys.Select) {
			m.DeleteLastIRQInput()
			return tea.Batch(cmds...)
		}

		// TODO: figure out why this isn't working (back button not working)
		if m.irqFocusIndex == backBtnIndex &&
			key.Matches(msg, m.keys.Select) {
			m.currMenu = mainMenu
			return tea.Batch(cmds...)
		}

		// Handle navigation between the buttons
		if m.irqFocusIndex == applyBtnIndex &&
			key.Matches(msg, m.keys.Left) {
			m.irqFocusIndex = backBtnIndex
		}
		if m.irqFocusIndex == backBtnIndex &&
			key.Matches(msg, m.keys.Up) {
			m.irqFocusIndex = plusBtnIndex
		}
		if m.irqFocusIndex == backBtnIndex &&
			key.Matches(msg, m.keys.Right) {
			m.irqFocusIndex = applyBtnIndex
		}
		if m.irqFocusIndex == plusBtnIndex &&
			key.Matches(msg, m.keys.Right) {
			m.irqFocusIndex = minusBtnIndex
		}
		if m.irqFocusIndex == minusBtnIndex &&
			key.Matches(msg, m.keys.Right) {
			m.irqFocusIndex = plusBtnIndex
		}

		// Validate the inputs
		// log.Println("IRQ menu isValid: ", m.AreValidInputs())

		// Handle [ Back ] button

		/* If the user press enter on the log view,
		go back to the previous menu */
		if m.renderLog && key.Matches(msg, m.keys.Select) {
			m.renderLog = false
			m.currMenu = kcmdlineMenu
		}

		// TODO: make field validation

		// Did the user press enter while the apply button was focused?
		// TODO: improve mapping of len(m.inputs) to the apply button
		if key.Matches(msg, m.keys.Select) &&
			m.irqFocusIndex == applyBtnIndex {

			// TODO: generate the IRQFilter structs based on the inputs
			log.Println("------Apply changes")

			// TODO: validation needs to be check here.

			// TODO: drop this
		}

		// Cycle indexes
		if key.Matches(msg, m.keys.Up) {
			// m.PrevIndex(&m.irqFocusIndex, m.irqInputs)
			prevFocusIndex(&m.irqFocusIndex, totalIrqItems)
		}

		if key.Matches(msg, m.keys.Down) ||
			key.Matches(msg, m.keys.Select) {
			// m.NextIndex(&m.irqFocusIndex, m.irqInputs)
			nextFoxcusIndex(&m.irqFocusIndex, totalIrqItems)
		}

		for i := 0; i <= len(m.irqInputs)-1; i++ {
			if i == m.irqFocusIndex {
				// Set focused state
				cmds[i] = m.irqInputs[i].Focus()
				m.irqInputs[i].PromptStyle = styles.FocusedStyle
				m.irqInputs[i].TextStyle = styles.FocusedStyle
				m.irqInputs[i].Placeholder = getPlaceholder(i)
				continue
			}
			// Remove focused state
			m.irqInputs[i].Blur()
			m.irqInputs[i].PromptStyle = styles.NoStyle
			m.irqInputs[i].TextStyle = styles.NoStyle
			m.irqInputs[i].Placeholder = ""
		}

	}
	for i := range m.irqInputs {
		m.irqInputs[i], cmds[i] = m.irqInputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func getPlaceholder(i int) string {
	if i%2 == 0 {
		return "Insert filter parameters for IRQs"
	}
	return cpuListPlaceholder
}
