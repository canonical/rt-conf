package ui

import (
	"log"

	cmp "github.com/canonical/rt-conf/src/ui/components"
	"github.com/canonical/rt-conf/src/ui/config"
	"github.com/canonical/rt-conf/src/ui/styles"
)

// TODO: Create a new object for the IRQinputs menu

// | IRQ Filter:
// | CPU Range:
// | [ - ]

// [ Back ] [ + ] [ Apply ]

// TODO: add [ - ] button to remove the last entry (or a specific entry)

// TODO: look into the pagination mechanism that bubbles provives for the list
// ** NOTE: It would be intesting to have a pagination mechanism for the
// ** text inputs

// TODO: delete this function
func (m *IRQMenuModel) EditModeView() string {
	log.Println("\n---- IRQEditMode VIEW ----")
	// var s string // the view

	// title := styles.InnerMenuStyle("Configuring IRQ Affinity")
	// desc := styles.Section.
	// 	Render("Allocate specific CPUs to IRQs matching the filter.")

	// // The inputs
	// var b strings.Builder
	// for i := range m.irqInputs {
	// 	// TODO fix the padding of the irqInputs view
	// 	textInput := m.irqInputs[i].View()
	// 	b.WriteString(textInput + "\n")
	// 	if (i+1)%2 == 0 {
	// 		b.WriteString("---\n")
	// 	}
	// }

	// plus_button := cmp.NewButton("+")
	// minus_button := cmp.NewButton("-")
	// apply_button := cmp.NewButton("Apply")
	// back_button := cmp.NewButton("Back")

	// btns := []*cmp.Button{plus_button, minus_button, back_button, apply_button}

	// for i, btn := range btns {
	// 	if m.irqFocusIndex == i+(len(m.irqInputs)) {
	// 		btn.SetFocused()
	// 	} else {
	// 		btn.SetBlurred()
	// 	}
	// }

	// // TODO: add space between the [ Apply ] and [ Back ] buttons
	// log.Println("(view) irqFocusIndex: ", m.irqFocusIndex)

	// fmt.Fprintf(&b, "\n%s\n",
	// 	styles.JoinHorizontal(
	// 		plus_button.Render(),
	// 		minus_button.Render(),
	// 	))

	// // [ Back ] [ Apply ] buttons
	// fmt.Fprintf(&b, "\n\n%s\n\n",
	// 	styles.JoinHorizontal(
	// 		back_button.Render(),
	// 		apply_button.Render(),
	// 	))

	// body := b.String()
	// // TODO: Adding padding to the bottom and top of [body] and remove new lines

	// helpView := m.help.View(m.keys)

	// // TODO: fix this mess
	// height := (m.height -
	// 	strings.Count(title, "\n") -
	// 	strings.Count(desc, "\n") -
	// 	strings.Count(helpView, "\n") -
	// 	strings.Count(m.errorMsg, "\n") -
	// 	strings.Count(body, "\n") - 6) / 2 // TODO: fix those magic numbers

	// // NOTE: *- 4 * because:
	// // "\n\n" (jumps 2 lines) after the title
	// // Before the line with the [ Back ] [ Apply ] buttons there are 2 lines

	// // NOTE: * / 2 * (divide by two) because:
	// // we want to add padding between to the top
	// // and bottom of the help view

	// if height < 0 {
	// 	height = 1
	// }

	// s +=
	// 	title +
	// 		"\n\n" +
	// 		desc +
	// 		"\n\n" +
	// 		body +
	// 		strings.Repeat("\n", height) +
	// 		"\n" +
	// 		// styles.ErrorMessageStyle(m.errorMsg) +
	// 		strings.Repeat("\n", height) +
	// 		helpView
	// return s
	return "**** UNDER CONSTRUCTION ****"
}

// TODO: Need to think a way to model the navigation between menus
func (m Model) View() string {
	log.Println("\n------------- VIEW -------------------")
	log.Println("(VIEW) Current stack: ", m.nav.PrintMenuStack())
	log.Printf("(VIEW) Current menu: %s", config.Menu[m.nav.GetCurrMenu()])

	switch m.nav.GetCurrMenu() {

	case config.KCMD_VIEW_ID:
		return styles.AppStyle.Render(m.kcmdlineView())

	case config.KCMD_CONCLUSSION_VIEW_ID:
		if m.renderLog {
			back_button := cmp.FocusedButton("Back")
			m.logMsg = append(m.logMsg, "\n")
			m.logMsg = append(m.logMsg, back_button)
			var content string
			for _, msg := range m.logMsg {
				content += msg
			}

			max := 0
			for _, msg := range m.logMsg {
				if len(msg) > max {
					max = len(msg)
				}
			}
			// Render the centered square with text
			return styles.CenteredSquareWithText(
				m.width, m.height, max, len(m.logMsg), content)
		}
		panic("\n\nInvalid view")

	case config.IRQ_VIEW_ID:
		log.Print("(view) Rendering IRQ Menu")
		return styles.AppStyle.Render(m.irq.list.View())

	case config.IRQ_ADD_EDIT_VIEW_ID:
		log.Print("(view) Rendering IRQ Add/Edit Menu")
		return styles.AppStyle.Render(m.irq.EditModeView())

	default: // INITIAL MENU
		return styles.AppStyle.Render(m.main.list.View())
	}
}
