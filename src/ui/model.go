package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	// WIP
}

func NewModel() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// No messages are currently handled.
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	return "WIP"
}
