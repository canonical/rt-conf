package ui

import (
	"fmt"
	"strings"

	cmp "github.com/canonical/rt-conf/src/ui/components"
	"github.com/canonical/rt-conf/src/ui/styles"
)

// TODO: look into the pagination mechanism that bubbles provives for the list
// ** NOTE: It would be intesting to have a pagination mechanism for the
// ** text inputs

func (m Model) irqTunningView() string {
	var s string // the view

	title := styles.InnerMenuStyle("Configuring IRQ Affinity")

	// The inputs
	var b strings.Builder
	for i := range m.irqInputs {
		// TODO fix the padding of the irqInputs view

		// TODO fix the violations of the index range
		// if i < len(m.irqInputs)-1 {
		// }
		if i < len(m.irqInputs)-1 {
			left := m.irqInputs[i].View()
			right := m.irqInputs[i+1].View()
			b.WriteString(styles.JoinHorizontal(&left, &right))
			b.WriteRune('\n')
		}
	}

	apply_button := cmp.NewButton("Apply")
	back_button := cmp.NewButton("Back")
	apply_button.SetBlurred()
	back_button.SetBlurred()

	// TODO: add space between the [ Apply ] and [ Back ] buttons
	if m.focusIndex == len(m.irqInputs) {
		apply_button.SetFocused()
		back_button.SetBlurred()

	} else if m.focusIndex == len(m.irqInputs)+1 {
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
	if m.focusIndex == len(m.kcmdInputs) {
		apply_button.SetFocused()
		back_button.SetBlurred()

	} else if m.focusIndex == len(m.kcmdInputs)+1 {
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
