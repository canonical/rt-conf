package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd { return nil }

func (m MainMenuModel) Init() tea.Cmd { return nil }

func (m KcmdlineMenuModel) Init() tea.Cmd { return textinput.Blink }

func (m KcmdlineConclussion) Init() tea.Cmd { return nil }

func (m IRQMenuModel) Init() tea.Cmd { return nil }

func (m IRQAddEditMenu) Init() tea.Cmd { return nil }
