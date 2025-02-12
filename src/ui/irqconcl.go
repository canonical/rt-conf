package ui

import (
	"log"
	"strconv"
	"strings"

	cmp "github.com/canonical/rt-conf/src/ui/components"
	"github.com/canonical/rt-conf/src/ui/styles"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type IRQConclusion struct {
	// The number of IRQs that are being applied to the system
	num    int
	Nav    *cmp.MenuNav // Menu Navigation instance
	keys   *irqKeyMap
	Width  int
	Height int
	logMsg string
	errMsg string
}

func newIRQConclusionModel() IRQConclusion {
	keys := irqMenuListKeyMap()
	nav := cmp.GetMenuNavInstance()
	return IRQConclusion{
		Nav:  nav,
		keys: keys,
	}
}

func (m IRQConclusion) Init() tea.Cmd { return nil }

func (m *IRQConclusion) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Select):
			m.Nav.PrevMenu()
		}
	}
	return m, nil
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
