package ui

// ** NOTE: this is being initializaed because the IRQinputs are always in pairs
// ** so the plusBtnIndex in the beggining will be: len(m.irqInputs) which is 2
// ** But this is necessary to avoid that the plus button is focused when the
// ** menu is rendered for the first time

// var plusBtnIndex = 2 // This is the initial index of the + button
// var minusBtnIndex = plusBtnIndex + 1
// var backBtnIndex = minusBtnIndex + 1
// var applyBtnIndex = backBtnIndex + 1

// ** NOTE this func may be useful
// // TODO: Functions to generate new entries for IRQ affinity menu
// func (m *Model) NewIRQTextInputs() {
// 	m.irqInputs = append(m.irqInputs, newIRQtextInputs()...)
// 	plusBtnIndex = len(m.irqInputs)
// 	minusBtnIndex = plusBtnIndex + 1
// 	m.irqFocusIndex = plusBtnIndex
// }

// ** NOTE this func may be useful
// func (m *Model) DeleteLastIRQInput() {
// 	if len(m.irqInputs) > 0 {
// 		m.irqInputs = m.irqInputs[:len(m.irqInputs)-2]
// 	}
// 	plusBtnIndex = len(m.irqInputs)
// 	minusBtnIndex = plusBtnIndex + 1
// 	m.irqFocusIndex = minusBtnIndex
// }

// // TODO: add exception to keys h,j,k,l not being handled as navigation keys
// func (m *Model) updateIRQMenu(msg tea.KeyMsg) tea.Cmd {
// 	log.Println("\n---- UPDATE ----")
// 	totalIrqItems := len(m.irqInputs) + 4
// 	index := components.NewNavigation(&m.irqFocusIndex, &totalIrqItems)
// 	cmds := make([]tea.Cmd, len(m.irqInputs))

// 	log.Println("Size of IRQ inputs: ", len(m.irqInputs))
// 	log.Println("(update) focusIndex: ", m.irqFocusIndex)

// 	plusBtnIndex = len(m.irqInputs)
// 	minusBtnIndex = plusBtnIndex + 1
// 	backBtnIndex = minusBtnIndex + 1
// 	applyBtnIndex = backBtnIndex + 1

// 	dbgHelper := make(map[int]string, len(m.irqInputs)+4)
// 	dbgHelper[plusBtnIndex] = "Plus btn"
// 	dbgHelper[minusBtnIndex] = "Minus btn"
// 	dbgHelper[backBtnIndex] = "Back btn"
// 	dbgHelper[applyBtnIndex] = "Apply btn"

// 	log.Printf("FocusIndex: %d = %s ", m.irqFocusIndex, dbgHelper[m.irqFocusIndex])

// 	switch {

// 	case key.Matches(msg, m.keys.Add):
// 		log.Println("------Add new IRQ input")
// 		// TODO: trigger logic to render the new IRQ inputs view

// 		// m.delegateKeys.remove.SetEnabled(true)
// 		// TODO: add logic to generate new IRQ fitler entry
// 		// newItem := m.itemGenerator.next()
// 		// insCmd := m.irq.list.InsertItem(0, newItem)
// 		// statusCmd := m.list.NewStatusMessage(styles.StatusMessageStyle("Added " + newItem.Title()))
// 		// return m, tea.Batch(insCmd, statusCmd)
// 		return tea.Batch(cmds...)

// 	case key.Matches(msg, m.keys.Apply):
// 		log.Println("------Apply changes")
// 		return tea.Batch(cmds...)

// 	case key.Matches(msg, m.keys.Select),
// 		key.Matches(msg, m.keys.Up),
// 		key.Matches(msg, m.keys.Down),
// 		key.Matches(msg, m.keys.Left),
// 		key.Matches(msg, m.keys.Right):

// 		// TODO: add logic for + btn
// 		// TODO: check weird behavior of this + button
// 		if m.irqFocusIndex == plusBtnIndex &&
// 			key.Matches(msg, m.keys.Select) {
// 			m.NewIRQTextInputs()
// 			return tea.Batch(cmds...)
// 		}
// 		if m.irqFocusIndex == minusBtnIndex &&
// 			key.Matches(msg, m.keys.Select) {
// 			m.DeleteLastIRQInput()
// 			return tea.Batch(cmds...)
// 		}

// 		// TODO: figure out why this isn't working (back button not working)
// 		if m.irqFocusIndex == backBtnIndex &&
// 			key.Matches(msg, m.keys.Select) {
// 			m.currMenu = mainMenu
// 			return tea.Batch(cmds...)
// 		}

// 		// Handle navigation between the buttons
// 		ret := m.handleBtnNav(msg, cmds, index)
// 		if ret != nil {
// 			return ret
// 		}

// 		// Validate the inputs
// 		// log.Println("IRQ menu isValid: ", m.AreValidInputs())

// 		// Handle [ Back ] button
// 		/* If the user press enter on the log view,
// 		go back to the previous menu */
// 		if m.renderLog && key.Matches(msg, m.keys.Select) {
// 			m.renderLog = false
// 			m.currMenu = kcmdlineMenu
// 		}

// 		// TODO: make field validation

// 		// Did the user press enter while the apply button was focused?
// 		// TODO: improve mapping of len(m.inputs) to the apply button
// 		if key.Matches(msg, m.keys.Select) &&
// 			m.irqFocusIndex == applyBtnIndex {
// 			// TODO: generate the IRQFilter structs based on the inputs
// 			log.Println("------Apply changes")

// 			// TODO: validation needs to be check here.
// 		}

// 		// Cycle indexes
// 		if key.Matches(msg, m.keys.Up) {
// 			// m.PrevIndex(&m.irqFocusIndex, m.irqInputs)
// 			log.Println("DIRECTION = UP")
// 			index.Prev()
// 		}

// 		if key.Matches(msg, m.keys.Down) ||
// 			key.Matches(msg, m.keys.Select) {
// 			// m.NextIndex(&m.irqFocusIndex, m.irqInputs)
// 			log.Println("DIRECTION = DOWN")
// 			index.Next()
// 		}

// 		for i := 0; i <= len(m.irqInputs)-1; i++ {
// 			if i == m.irqFocusIndex {
// 				// Set focused state
// 				cmds[i] = m.irqInputs[i].Focus()
// 				m.irqInputs[i].PromptStyle = styles.FocusedStyle
// 				m.irqInputs[i].TextStyle = styles.FocusedStyle
// 				m.irqInputs[i].Placeholder = getPlaceholder(i)
// 				continue
// 			}
// 			// Remove focused state
// 			m.irqInputs[i].Blur()
// 			m.irqInputs[i].PromptStyle = styles.NoStyle
// 			m.irqInputs[i].TextStyle = styles.NoStyle
// 			m.irqInputs[i].Placeholder = ""
// 		}

// 	}
// 	for i := range m.irqInputs {
// 		m.irqInputs[i], cmds[i] = m.irqInputs[i].Update(msg)
// 	}
// 	return tea.Batch(cmds...)
// }

// func getPlaceholder(i int) string {
// 	if i%2 == 0 {
// 		return "Insert filter parameters for IRQs"
// 	}
// 	return cpuListPlaceholder
// }

// func (m *Model) handleBtnNav(msg tea.KeyMsg, cmds []tea.Cmd,
// 	index *components.Navigation) tea.Cmd {
// 	// Handle navigation between the buttons
// 	if m.irqFocusIndex == applyBtnIndex &&
// 		key.Matches(msg, m.keys.Left) {
// 		index.Prev()
// 		return tea.Batch(cmds...)
// 	}
// 	if m.irqFocusIndex == applyBtnIndex &&
// 		key.Matches(msg, m.keys.Right) {
// 		index.Next()
// 		return tea.Batch(cmds...)
// 	}
// 	if m.irqFocusIndex == backBtnIndex &&
// 		key.Matches(msg, m.keys.Right) {
// 		index.Next()
// 		return tea.Batch(cmds...)
// 	}
// 	if m.irqFocusIndex == backBtnIndex &&
// 		key.Matches(msg, m.keys.Left) {
// 		index.Prev()
// 		return tea.Batch(cmds...)
// 	}
// 	if m.irqFocusIndex == plusBtnIndex &&
// 		key.Matches(msg, m.keys.Right) {
// 		index.Next()
// 		return tea.Batch(cmds...)
// 	}
// 	if m.irqFocusIndex == plusBtnIndex &&
// 		key.Matches(msg, m.keys.Left) {
// 		index.Prev()
// 		return tea.Batch(cmds...)
// 	}
// 	if m.irqFocusIndex == minusBtnIndex &&
// 		key.Matches(msg, m.keys.Right) {
// 		index.Next()
// 		return tea.Batch(cmds...)
// 	}
// 	if m.irqFocusIndex == minusBtnIndex &&
// 		key.Matches(msg, m.keys.Left) {
// 		index.Prev()
// 		return tea.Batch(cmds...)
// 	}
// 	return nil
// }

// func (m Model) UpdateIRQMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	return m, nil
// }

// func (m Model) UpdateIRQMenu(msg tea.Msg) tea.Cmd {
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		if msg.String() == "ctrl+c" {
// 			return tea.Quit
// 		}
// 	}

// 	var cmd tea.Cmd
// 	m.irqMenu, cmd = m.irqMenu.Update(msg)
// 	return cmd
// }
