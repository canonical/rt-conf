package model

import (
	"fmt"
	"strconv"
	"strings"

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
	ScalGov string `yaml:"scaling-governor"`
	MinFreq string `yaml:"min-freq"`
	MaxFreq string `yaml:"max-freq"`
}

func (c CpuGovernanceRule) Validate() error {
	if _, err := cpulists.Parse(c.CPUs); err != nil {
		return err
	}

	if _, ok := scalProfilesMap[c.ScalGov]; !ok && c.ScalGov != "" {
		return fmt.Errorf("invalid cpu scaling governor: %v", c.ScalGov)
	}

	minFreq, err := ParseFreq(c.MinFreq)
	if err != nil {
		return fmt.Errorf("invalid min frequency: %v", err)
	}
	maxFreq, err := ParseFreq(c.MaxFreq)
	if err != nil {
		return fmt.Errorf("invalid max frequency: %v", err)
	}

	if err := validateFreqRange(minFreq, maxFreq); err != nil {
		return fmt.Errorf("invalid frequency range: %v", err)
	}

	return nil
}

func validateFreqRange(min, max int) error {
	if min == -1 && max == -1 {
		return nil // No frequency bounds
	}

	if (min != -1 && max != -1) && max < min {
		return fmt.Errorf(
			"max frequency (%d) should not be less than min frequency (%d)",
			max, min)
	}

	return nil
}

func ParseFreq(freq string) (int, error) {
	if freq == "" {
		return -1, nil // No frequency limits set, nothing to parse
	}

	s := strings.ToLower(strings.TrimSpace(freq))
	if !strings.HasSuffix(s, "hz") {
		return -1, fmt.Errorf(
			"invalid format: frequency must end with 'Hz': %s", s)
	}

	s = strings.TrimSuffix(s, "hz")

	multiplier := 1.0
	switch {
	case strings.HasSuffix(s, "g"):
		multiplier = 1_000_000.0
		s = strings.TrimSuffix(s, "g")
	case strings.HasSuffix(s, "m"):
		multiplier = 1_000.0
		s = strings.TrimSuffix(s, "m")
	case strings.HasSuffix(s, "k"):
		multiplier = 1.0
		s = strings.TrimSuffix(s, "k")
	default:
		multiplier = 0.001 // Default to raw Hz if no metric prefix is provided
	}

	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return -1, fmt.Errorf("failed to parse frequency value: %v", err)
	}

	if val < 0 {
		return -1, fmt.Errorf("frequency value cannot be negative: %s", s)
	}

	kHz := int(val * multiplier)
	return kHz, nil
}
