package ui

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

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
