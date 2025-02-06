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

// TODO: URGENT = [ BACK ] isn't working
func (m Model) kcmdlineView() string {
	var s string // the view

	title := styles.InnerMenuStyle("Configuring Kernel Cmdline Parameters")

	// The inputs
	var b strings.Builder
	for i := range m.kcmd.Inputs {
		b.WriteString(m.kcmd.Inputs[i].View())
		if i < len(m.kcmd.Inputs)-1 {
			b.WriteRune('\n')
		}
	}

	apply_button := cmp.NewButton("Apply")
	back_button := cmp.NewButton("Back")
	apply_button.SetBlurred()
	back_button.SetBlurred()

	// TODO: add space between the [ Apply ] and [ Back ] buttons
	if m.kcmd.FocusIndex == len(m.kcmd.Inputs) {
		apply_button.SetFocused()
		back_button.SetBlurred()

	} else if m.kcmd.FocusIndex == len(m.kcmd.Inputs)+1 {
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

	//TODO: this needs to be dropped
	helpView := m.kcmd.help.View(m.kcmd.keys)

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
