package data

import (
	"fmt"

	"github.com/canonical/rt-conf/src/cpu"
)

type ScalProfiles int

const (
	balanced ScalProfiles = iota
	powersave
	performance
)

var ScalProfilesMap = map[string]ScalProfiles{
	"balanced":    balanced,
	"powersave":   powersave,
	"performance": performance,
}

type CpuGovernance struct {
	CPUs    string `yaml:"cpus"`
	ScalGov string `yaml:"scaling_governor"`
}

func (c CpuGovernance) Validate() error {
	if _, ok := ScalProfilesMap[c.ScalGov]; !ok {
		return fmt.Errorf("invalid cpu scaling governor: %v", c.ScalGov)
	}
	err := cpu.ValidateList(c.CPUs)
	if err != nil {
		return fmt.Errorf("invalid cpus: %v", err)
	}
	return nil
}
