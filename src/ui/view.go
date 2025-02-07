package ui

import (
	"fmt"
	"log"
	"strings"

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
	var s string // the view

	title := styles.InnerMenuStyle("Configuring IRQ Affinity")
	desc := styles.Section.
		Render("Allocate specific CPUs to IRQs matching the filter.")

	var b strings.Builder
	for i := range m.irq.Inputs {
		b.WriteString(m.irq.Inputs[i].View())
		if i < len(m.irq.Inputs)-1 {
			b.WriteRune('\n')
		}
	}
	b.WriteString("\n")

	var verticalPadding strings.Builder
	verticalPadding.WriteString("\n\n")

	okayBtn := cmp.NewButton("Okay")
	cancelBtn := cmp.NewButton("Cancel")

	btns := []*cmp.Button{okayBtn, cancelBtn}

	for i, btn := range btns {
		if m.irq.FocusIndex == i+(len(m.irq.Inputs)) {
			btn.SetFocused()
		} else {
			btn.SetBlurred()
		}
	}

	log.Println("--- (IRQ ADD/EDIT VIEW) m.irq.FocusIndex: ", m.irq.FocusIndex)

	fmt.Fprintf(&b, "\n%s\n",
		styles.JoinHorizontal(
			okayBtn.Render(),
			cancelBtn.Render(),
		))

	body := b.String()

	helpView := m.help.View(m.keys)

	height := m.height -
		strings.Count(title, "\n") -
		strings.Count(desc, "\n") -
		strings.Count(helpView, "\n") -
		strings.Count(m.irq.errorMsgCpu, "\n") -
		strings.Count(m.irq.errorMsgFilter, "\n") -
		strings.Count(b.String(), "\n") -
		// verticalPadding is used twice
		strings.Count(verticalPadding.String(), "\n") -
		strings.Count(verticalPadding.String(), "\n")

	log.Println("--- (IRQ ADD/EDIT VIEW) m.height: ", m.height)
	log.Println("--- (IRQ ADD/EDIT VIEW) height: ", height)

	if height < 0 {
		height = 1
	}

	s +=
		title +
			verticalPadding.String() +
			desc +
			verticalPadding.String() +
			body +
			strings.Repeat("\n", height/2) +
			styles.ErrorMessageStyle(m.irq.errorMsgFilter) +
			styles.ErrorMessageStyle(m.irq.errorMsgCpu) +
			strings.Repeat("\n", height/2) +
			helpView
	return s
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
