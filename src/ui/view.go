package ui

import (
	"fmt"
	"strings"

	cmp "github.com/canonical/rt-conf/src/ui/components"
	"github.com/canonical/rt-conf/src/ui/styles"
)

// TODO: fix the padding of the kcmline view
// * NOTE: in comparison with the main menu, the title is shifted to left
// * Not only the tittle but the hole view is shifted to the left

func (m Model) kcmdlineView() string {
	var s string // the view
	// m.infoMsg = "\n"

	title := styles.InnerMenuStyle("Configuring Kernel Cmdline Parameters")

	// The inputs
	var b strings.Builder
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	apply_button := cmp.NewButton("Apply")
	back_button := cmp.NewButton("Back")
	apply_button.SetBlurred()
	back_button.SetBlurred()

	// TODO: add space between the [ Apply ] and [ Back ] buttons
	if m.focusIndex == len(m.inputs) {
		apply_button.SetFocused()
		back_button.SetBlurred()

	} else if m.focusIndex == len(m.inputs)+1 {
		apply_button.SetBlurred()
		back_button.SetFocused()
	}

	// [ Back ] [ Apply ] buttons
	fmt.Fprintf(&b, "\n\n%s\n\n",
		styles.JoinHorizontal(
			back_button.Render(),
			apply_button.Render(),
		))

	// Cursor mode message //TODO: MAYBE use this to display info
	// b.WriteString(styles.HelpStyle.Render("cursor mode is "))
	// b.WriteString(styles.CursorModeHelpStyle.Render(m.cursorMode.String()))
	// b.WriteString(styles.HelpStyle.Render(" (ctrl+r to change style)"))

	body := b.String()
	// TODO: Adding padding to the bottom and top of [body] and remove new lines

	// The help view
	helpView := m.help.View(m.keys)

	// TODO: fix this mess
	height := (m.height -
		strings.Count(title, "\n") -
		strings.Count(helpView, "\n") -
		strings.Count(m.errorMsg, "\n") -
		strings.Count(body, "\n") - 6) / 2 // TODO: fix those magic numbers

	// NOTE: *- 6 * because:
	// "\n\n" (jumps 2 lines)

	// there is the line with the [ Back ] [ Apply ] buttons
	// between the m.infoMsg and m.errorMsg there is 1 line
	// between the m.logMsg and m.infoMsg there is 1 line

	// NOTE: *- 8 * because:
	// There is 8 lines of output when processing the kcmdline functions

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
			// logMessageStyle(m.logMsg) +
			"\n" +
			// styles.InfoMessageStyle(m.infoMsg) +
			"\n" +
			"\n" +
			styles.ErrorMessageStyle(m.errorMsg) +
			strings.Repeat("\n", height) +
			helpView
	return s
}

func (m Model) irqAffinityView() string {

	title := styles.InnerMenuStyle("Configuring IRQ Affinity")

	helpView := m.help.View(m.keys)
	height := m.height - strings.Count(title, "\n") - strings.Count(helpView, "\n")

	return "\n" + title + strings.Repeat("\n", height) + helpView
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
		return styles.AppStyle.Render(m.irqAffinityView())
	default:
		return styles.AppStyle.Render(m.list.View())
	}
}
