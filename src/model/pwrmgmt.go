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
	MinFreq string `yaml:"min_freq"`
	MaxFreq string `yaml:"max_freq"`
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
	if freq == "" {
		return nil // No frequency limits set, nothing to validate
	}
	reg := regexp.MustCompile(`^\d*\.?\d*[KkGgMm]{1}([Hh][Zz]){1}$`)
	if !reg.MatchString(freq) {
		msg := "expected formats: 3.4GHz, 2000MHz, 100000KHz, got: " + freq
		return fmt.Errorf("invalid frequency format: %s", msg)
	}
	return nil
}
