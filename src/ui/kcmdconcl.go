package ui

import (
	"log"

	cmp "github.com/canonical/rt-conf/src/ui/components"
	"github.com/canonical/rt-conf/src/ui/styles"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type KcmdlineConclusion struct {
	Nav       *cmp.MenuNav // Menu Navigation instance
	keys      *kcmdKeyMap
	Width     int
	Height    int
	logMsg    []string
	renderLog bool
}

func newKcmdConclusionModel() KcmdlineConclusion {
	keys := newkcmdMenuListKeyMap()
	nav := cmp.GetMenuNavInstance()
	return KcmdlineConclusion{
		Nav:  nav,
		keys: keys,
	}
}

func (m KcmdlineConclusion) Init() tea.Cmd { return nil }

func (m *KcmdlineConclusion) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	log.Println("(kcmdlineConclusionUpdate - start")

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		/* If the user press enter on the log view,
		go back to the previous menu */
		case key.Matches(msg, m.keys.Select):
			m.Nav.PrevMenu()

		default:
			log.Println("(default) Invalid key pressed: ", msg.String())
			log.Println("Menu navigation stack: ", m.Nav.PrintMenuStack())
		}
	}
	return m, cmd
}

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
