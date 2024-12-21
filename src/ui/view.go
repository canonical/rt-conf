package ui

import (
	"fmt"
	"strings"

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

	title := innerMenuStyle("Configuring Kernel Cmdline Parameters")

	// The inputs
	var b strings.Builder
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)
	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))
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
			infoMessageStyle(m.infoMsg) +
			"\n" +
			errorMessageStyle(m.errorMsg) +
			strings.Repeat("\n", height) +
			bottom
	return s
}

func (m Model) irqAffinityView() string {

	title := innerMenuStyle("Configuring IRQ Affinity")

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
			return CenteredSquareWithText(
				m.width, m.height, max, len(m.logMsg), content)
		}
		// TODO: check for the [ apply ] button then show a clean view
		// to show the needed actions from the user
		return appStyle.Render(m.kcmdlineView())
	case irqAffinityMenu:
		return appStyle.Render(m.irqAffinityView())
	default:
		return appStyle.Render(m.list.View())
	}
}
