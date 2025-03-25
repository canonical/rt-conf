package ui

import (
	"fmt"
	"log"

	"github.com/canonical/rt-conf/src/model"
	tea "github.com/charmbracelet/bubbletea"
)

func Start(conf *model.InternalConfig) error {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		return fmt.Errorf("failed to open log file: %s", err)
	}
	defer f.Close()

	// Run the Terminal User Interface (TUI)
	// This is a blocking call
	if _, err := tea.NewProgram(NewModel(conf), tea.WithAltScreen()).Run(); err != nil {
		log.Fatalf("rt-conf failed: %v", err)
		return fmt.Errorf("error creating the UI program: %s", err)
	}

	return nil
}
