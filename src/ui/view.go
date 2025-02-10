package ui

import (
	"log"
	"strconv"
	"strings"

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

func (m IRQConclusion) View() string {
	// TODO: fix this view
	log.Println("---- IRQConclusion VIEW ----")

	var s string
	backBtn := cmp.FocusedButton("Back")

	log.Println("Number of rules added: ", m.num)
	s +=
		m.logMsg +
			strconv.Itoa(m.num) +
			"\n" +
			m.errMsg +
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
