package ui

import (
	"log"

	"github.com/canonical/rt-conf/src/ui/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// TODO: add conclussion screen saying:
// "5 IRQ filter rules are aplied to the system"
// [ BACK ]

type NewIRQEntryMsg struct {
	irqAffinityRule
}

type irqAffinityRule struct {
	filter, cpulist string
}

func (i irqAffinityRule) Title() string       { return i.filter }
func (i irqAffinityRule) Description() string { return i.cpulist }

// This needs to be implemented to satisfy the list.Item interface
func (i irqAffinityRule) FilterValue() string { return i.filter }


