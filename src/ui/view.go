package ui

import (
	cmp "github.com/canonical/rt-conf/src/ui/components"
	"github.com/canonical/rt-conf/src/ui/styles"
)

func (m KcmdlineConclusion) View() string {
	if !m.renderLog {
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
