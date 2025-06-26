package model

import (
	"fmt"
	"regexp"

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
	MinFreq string `yaml:"min_frequency"`
	MaxFreq string `yaml:"max_frequency"`
}

func (c CpuGovernanceRule) Validate() error {
	if _, ok := scalProfilesMap[c.ScalGov]; !ok {
		return fmt.Errorf("invalid cpu scaling governor: %v", c.ScalGov)
	}
	_, err := cpulists.Parse(c.CPUs)
	if err != nil {
		return fmt.Errorf("invalid cpus: %v", err)
	}
	if err := c.CheckFreqFormat(); err != nil {
		return err
	}
	return nil
}

func (c CpuGovernanceRule) CheckFreqFormat() error {
	if err := CheckFreqFormat(c.MinFreq); err != nil {
		return fmt.Errorf("invalid min frequency: %w", err)
	}
	if err := CheckFreqFormat(c.MaxFreq); err != nil {
		return fmt.Errorf("invalid max frequency: %w", err)
	}
	return nil
}

func CheckFreqFormat(freq string) error {
	reg := regexp.MustCompile(`^\d*\.?\d*[GgMm]?([Hh][Zz])?$`)
	if !reg.MatchString(freq) {
		return fmt.Errorf("invalid frequency format: %s", freq)
	}
	return nil
}
