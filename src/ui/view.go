package ui

import (
	"fmt"
	"strings"

	"github.com/canonical/rt-conf/src/ui/styles"
	"github.com/charmbracelet/bubbles/key"
)

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k listKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.toggleHelpMenu, k.selectMenu}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k listKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.goHome, k.selectMenu}, // first column
		{k.toggleHelpMenu},       // second column
	}
}

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

	button := &styles.BlurredButton
	if m.focusIndex == len(m.inputs) {
		button = &styles.FocusedButton
	}

	// [ Apply ] button
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	// Cursor mode message
	b.WriteString(styles.HelpStyle.Render("cursor mode is "))
	b.WriteString(styles.CursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(styles.HelpStyle.Render(" (ctrl+r to change style)"))
	body := b.String()
	// TODO: Adding padding to the bottom and top of [body] and remove new lines

	// The help view
	helpView := m.help.View(m.keys)

	height := (m.height -
		strings.Count(title, "\n") -
		strings.Count(helpView, "\n") -
		// strings.Count(m.logMsg, "\n") -
		strings.Count(m.infoMsg, "\n") -
		strings.Count(m.errorMsg, "\n") -
		strings.Count(body, "\n") - 5) / 2 // TODO: fix those magic numbers

	// NOTE: *- 5 * because:
	// "\n\n" (jumps 2 lines)
	// m.infoMsg is always 1 line
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

	bottom := helpView

	s +=
		title +
			"\n\n" +
			body +
			strings.Repeat("\n", height) +
			// logMessageStyle(m.logMsg) +
			"\n" +
			styles.InfoMessageStyle(m.infoMsg) +
			"\n" +
			styles.ErrorMessageStyle(m.errorMsg) +
			strings.Repeat("\n", height) +
			bottom
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
		if m.renderLog {

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
		// TODO: check for the [ apply ] button then show a clean view
		// to show the needed actions from the user
		return styles.AppStyle.Render(m.kcmdlineView())
	case irqAffinityMenu:
		return styles.AppStyle.Render(m.irqAffinityView())
	default:
		return styles.AppStyle.Render(m.list.View())
	}
}
