package model

import (
	"fmt"

	"github.com/canonical/rt-conf/src/cpulists"
)

type ScalProfiles int

const (
	balanced ScalProfiles = iota
	powersave
	performance
)

var scalProfilesMap = map[string]ScalProfiles{
	"balanced":    balanced,
	"powersave":   powersave,
	"performance": performance,
}

type CpuGovernanceRule struct {
	CPUs    string `yaml:"cpus"`
	ScalGov string `yaml:"scaling_governor"`
	MaxFreq string `yaml:"max_frequency"`
	MinFreq string `yaml:"min_frequency"`
}

func (c CpuGovernanceRule) Validate() error {
	if _, ok := scalProfilesMap[c.ScalGov]; !ok {
		return fmt.Errorf("invalid cpu scaling governor: %v", c.ScalGov)
	}
	_, err := cpulists.Parse(c.CPUs)
	if err != nil {
		return fmt.Errorf("invalid cpus: %v", err)
	}
	return nil
}
