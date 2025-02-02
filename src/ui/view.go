package ui

import (
	"fmt"
	"log"
	"strings"

	cmp "github.com/canonical/rt-conf/src/ui/components"
	"github.com/canonical/rt-conf/src/ui/styles"
)

// TODO: add [ - ] button to remove the last entry (or a specific entry)

// TODO: look into the pagination mechanism that bubbles provives for the list
// ** NOTE: It would be intesting to have a pagination mechanism for the
// ** text inputs

func (m Model) irqTunningView() string {
	var s string // the view

	title := styles.InnerMenuStyle("Configuring IRQ Affinity")
	desc := styles.Section.
		Render("Allocate specific CPUs to IRQs matching the filter.")

	// The inputs
	var b strings.Builder
	for i := range m.irqInputs {
		// TODO fix the padding of the irqInputs view
		textInput := m.irqInputs[i].View()
		b.WriteString(textInput + "\n")
		if (i+1)%2 == 0 {
			b.WriteString("---\n")
		}
	}

	plus_button := cmp.NewButton("+")
	minus_button := cmp.NewButton("-")
	apply_button := cmp.NewButton("Apply")
	back_button := cmp.NewButton("Back")

	btns := []*cmp.Button{plus_button, minus_button, apply_button, back_button}

	// plus_button.SetBlurred()
	// apply_button.SetBlurred()
	// back_button.SetBlurred()
	// minus_button.SetBlurred()

	for _, btn := range btns {
		btn.SetBlurred()
	}

	for i, btn := range btns {
		if m.irqFocusIndex == i+(len(m.irqInputs)) {
			btn.SetFocused()
		} else {
			btn.SetBlurred()
		}
	}

	// TODO: add space between the [ Apply ] and [ Back ] buttons
	log.Println("irqFocusIndex: ", m.irqFocusIndex)

	// switch {
	// case m.irqFocusIndex == plusBtnIndex:
	// 	plus_button.SetFocused()
	// 	apply_button.SetBlurred()
	// 	back_button.SetBlurred()
	// 	minus_button.SetBlurred()
	// case m.irqFocusIndex == minusBtnIndex:
	// 	minus_button.SetFocused()
	// 	plus_button.SetBlurred()
	// 	apply_button.SetBlurred()
	// 	back_button.SetBlurred()
	// case m.irqFocusIndex == applyBtnIndex:
	// 	apply_button.SetFocused()
	// 	back_button.SetBlurred()
	// 	plus_button.SetBlurred()
	// 	minus_button.SetBlurred()
	// case m.irqFocusIndex == backBtnIndex:
	// 	back_button.SetFocused()
	// 	apply_button.SetBlurred()
	// 	plus_button.SetBlurred()
	// 	minus_button.SetBlurred()
	// }

	// [ + ] button
	// fmt.Fprintf(&b, "\n%s\n", *plus_button.Render())

	fmt.Fprintf(&b, "\n%s\n",
		styles.JoinHorizontal(
			plus_button.Render(),
			minus_button.Render(),
		))

	// [ Back ] [ Apply ] buttons
	fmt.Fprintf(&b, "\n\n%s\n\n",
		styles.JoinHorizontal(
			back_button.Render(),
			apply_button.Render(),
		))

	body := b.String()
	// TODO: Adding padding to the bottom and top of [body] and remove new lines

	helpView := m.help.View(m.keys)

	// TODO: fix this mess
	height := (m.height -
		strings.Count(title, "\n") -
		strings.Count(desc, "\n") -
		strings.Count(helpView, "\n") -
		strings.Count(m.errorMsg, "\n") -
		strings.Count(body, "\n") - 6) / 2 // TODO: fix those magic numbers

	// NOTE: *- 4 * because:
	// "\n\n" (jumps 2 lines) after the title
	// Before the line with the [ Back ] [ Apply ] buttons there are 2 lines

	// NOTE: * / 2 * (divide by two) because:
	// we want to add padding between to the top
	// and bottom of the help view

	if height < 0 {
		height = 1
	}

	s +=
		title +
			"\n\n" +
			desc +
			"\n\n" +
			body +
			strings.Repeat("\n", height) +
			"\n" +
			styles.ErrorMessageStyle(m.errorMsg) +
			strings.Repeat("\n", height) +
			helpView
	return s
}

// TODO: fix the padding of the kcmline view
// * NOTE: in comparison with the main menu, the title is shifted to left
// * Not only the tittle but the hole view is shifted to the left

func (m Model) kcmdlineView() string {
	var s string // the view

	title := styles.InnerMenuStyle("Configuring Kernel Cmdline Parameters")

	// The inputs
	var b strings.Builder
	for i := range m.kcmdInputs {
		b.WriteString(m.kcmdInputs[i].View())
		if i < len(m.kcmdInputs)-1 {
			b.WriteRune('\n')
		}
	}

	apply_button := cmp.NewButton("Apply")
	back_button := cmp.NewButton("Back")
	apply_button.SetBlurred()
	back_button.SetBlurred()

	// TODO: add space between the [ Apply ] and [ Back ] buttons
	if m.kcmdFocusIndex == len(m.kcmdInputs) {
		apply_button.SetFocused()
		back_button.SetBlurred()

	} else if m.kcmdFocusIndex == len(m.kcmdInputs)+1 {
		apply_button.SetBlurred()
		back_button.SetFocused()
	}

	// [ Back ] [ Apply ] buttons
	fmt.Fprintf(&b, "\n\n%s\n\n",
		styles.JoinHorizontal(
			back_button.Render(),
			apply_button.Render(),
		))

	body := b.String()
	// TODO: Adding padding to the bottom and top of [body] and remove new lines

	helpView := m.help.View(m.keys)

	// TODO: fix this mess
	height := (m.height -
		strings.Count(title, "\n") -
		strings.Count(helpView, "\n") -
		strings.Count(m.errorMsg, "\n") -
		strings.Count(body, "\n") - 4) / 2 // TODO: fix those magic numbers

	// NOTE: *- 4 * because:
	// "\n\n" (jumps 2 lines) after the title
	// Before the line with the [ Back ] [ Apply ] buttons there are 2 lines

	// NOTE: * / 2 * (divide by two) because:
	// we want to add padding between to the top
	// and bottom of the help view

	if height < 0 {
		height = 1
	}

	s +=
		title +
			"\n\n" +
			body +
			strings.Repeat("\n", height) +
			"\n" +
			styles.ErrorMessageStyle(m.errorMsg) +
			strings.Repeat("\n", height) +
			helpView
	return s
}

// TODO: Need to think a way to model the navigation between menus
func (m Model) View() string {
	switch m.currMenu {
	case kcmdlineMenu:

		// TODO: move this to a separate function
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
		return styles.AppStyle.Render(m.kcmdlineView())
	case irqAffinityMenu:
		return styles.AppStyle.Render(m.irqTunningView())
	default:
		return styles.AppStyle.Render(m.list.View())
	}
}
