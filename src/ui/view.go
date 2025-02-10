package ui

import (
	"fmt"
	"log"
	"strconv"
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
func (m IRQAddEditMenu) View() string {
	log.Println("\n---- IRQEditMode VIEW ----")
	var s string // the view

	title := styles.InnerMenuStyle("Add IRQ Affinity Rule")
	desc := styles.Section.
		Render("Allocate specific CPUs to IRQs matching the filter.")

	// TODO: bold cpu range and filter
	var b strings.Builder
	for i := range m.Inputs {
		b.WriteString(m.Inputs[i].View())
		if i < len(m.Inputs)-1 {
			b.WriteRune('\n')
		}
	}
	b.WriteString("\n")

	var verticalPadding strings.Builder
	verticalPadding.WriteString("\n\n")

	addBtn := cmp.NewButton("Add")
	cancelBtn := cmp.NewButton("Cancel")

	btns := []*cmp.Button{addBtn, cancelBtn}

	for i, btn := range btns {
		if m.FocusIndex == i+(len(m.Inputs)) {
			btn.SetFocused()
		} else {
			btn.SetBlurred()
		}
	}

	// log.Println("--- (IRQ ADD/EDIT VIEW) m.irq.FocusIndex: ", m.FocusIndex)

	fmt.Fprintf(&b, "\n%s\n",
		styles.JoinHorizontal(
			addBtn.Render(),
			cancelBtn.Render(),
		))

	body := b.String()

	helpView := m.help.View(m.keys)

	height := m.height -
		strings.Count(title, "\n") -
		strings.Count(desc, "\n") -
		strings.Count(helpView, "\n") -
		strings.Count(m.errorMsgCpu, "\n") -
		strings.Count(m.errorMsgFilter, "\n") -
		strings.Count(b.String(), "\n") -
		// verticalPadding is used twice
		strings.Count(verticalPadding.String(), "\n") -
		strings.Count(verticalPadding.String(), "\n")

	// log.Println("--- (IRQ ADD/EDIT VIEW) m.Height: ", m.height)
	// log.Println("--- (IRQ ADD/EDIT VIEW) height: ", height)

	if height < 0 {
		height = 1
	}
	// log.Println("--- recalculated height: ", height)

	s +=
		title +
			verticalPadding.String() +
			desc +
			verticalPadding.String() +
			body +
			strings.Repeat("\n", height/2) +
			styles.ErrorMessageStyle(m.errorMsgFilter) +
			styles.ErrorMessageStyle(m.errorMsgCpu) +
			strings.Repeat("\n", height/2) +
			helpView

	return styles.AppStyle.Render(s)
}

func (m MainMenuModel) View() string {
	return styles.AppStyle.Render(m.list.View())
}

func (m Model) View() string {
	// log.Println("\n------------- VIEW -------------------")
	// log.Printf("(VIEW) Current menu: %s", config.Menu[m.Nav.GetCurrMenu()])
	return m.GetActiveMenu().View()
}

func (m KcmdlineConclussion) View() string {
	if m.renderLog {
		panic("This call is not expected")
	}

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
		m.Width, m.Height, max, len(m.logMsg), content)
}

func (m IRQMenuModel) View() string {
	m.list.KeyMap.Filter.SetEnabled(false)
	m.list.KeyMap.CursorUp.SetEnabled(false)
	m.list.KeyMap.CursorDown.SetEnabled(false)
	m.list.KeyMap.NextPage.SetEnabled(false)
	m.list.KeyMap.PrevPage.SetEnabled(false)
	m.list.KeyMap.GoToStart.SetEnabled(false)
	m.list.KeyMap.GoToEnd.SetEnabled(false)
	m.list.KeyMap.Filter.SetEnabled(false)
	return styles.AppStyle.Render(m.list.View())
}

func (m IRQConclussion) View() string {
	// TODO: fix this view
	log.Println("---- IRQConclussion VIEW ----")

	var s string
	backBtn := cmp.FocusedButton("Back")

	log.Println("Number of rules added: ", m.num)
	s +=
		strconv.Itoa(m.num) +
			" " +
			m.logMsg +
			"\n\n" +
			backBtn

	log.Println("Size of s string: ", len(s))
	log.Println("Size of logMsg string: ", len(m.logMsg))

	lines := strings.Split(s, "\n")
	maxLenght := 0
	for _, line := range lines {
		if l := len([]rune(line)); l > maxLenght {
			maxLenght = l
		}
	}

	hight := strings.Count(s, "\n")
	// ** Note this -1 is a workaround to avoid
	// **  the top line margin to be cut off
	return styles.CenteredSquareWithText(
		m.Width, m.Height-1, maxLenght, hight, s)
}
