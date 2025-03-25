package ui

import (
	"github.com/canonical/rt-conf/src/model"
)

type IRQRuleMsg struct {
	IRQAffinityRule
}

type IRQAffinityRule struct {
	edited          bool // false means it's a new rule, true means new rule
	index           int
	rule            model.IRQTuning
	filter, cpulist string
}

type StartNewIRQAffinityRule struct {
	editMode bool
	// ** If editMode is true, this is the index of the rule
	// **  to edit as well the cpu list and filter
	index           int
	filter, cpulist string
}

func (i IRQAffinityRule) Title() string       { return i.filter }
func (i IRQAffinityRule) Description() string { return i.cpulist }

// This needs to be implemented to satisfy the list.Item interface
func (i IRQAffinityRule) FilterValue() string { return i.filter }
