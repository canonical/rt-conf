package model

import (
	"fmt"
	"reflect"
	"regexp"
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
	CPUs    string `yaml:"cpus" validation:"cpulist"`
	ScalGov string `yaml:"scaling_governor" validation:"scalgov"`
	MinFreq string `yaml:"min_freq" validation:"freq"`
	MaxFreq string `yaml:"max_freq" validation:"freq"`
}

func CheckFreqFormat(freq string) error {
	if freq == "" {
		return nil // No frequency limits set, nothing to validate
	}
	reg := regexp.MustCompile(`^\d+\.?\d*[KkGgMm]{1}([Hh][Zz]){1}$`)
	if !reg.MatchString(freq) {
		msg := "expected formats: 3.4GHz, 2000MHz, 100000KHz, got: " + freq
		return fmt.Errorf("invalid frequency format: %s", msg)
	}
	return nil
}

func (c CpuGovernanceRule) Validate() error {
	v := reflect.ValueOf(c)
	t := reflect.TypeOf(c)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		val := v.Field(i).String()
		tag := field.Tag.Get("validation")

		switch tag {
		case "cpulist":
			if _, err := cpulists.Parse(val); err != nil {
				return err
			}
		case "scalgov":
			if _, ok := scalProfilesMap[val]; !ok && val != "" {
				return fmt.Errorf("invalid cpu scaling governor: %v", val)
			}
		case "freq":
			if err := CheckFreqFormat(val); err != nil {
				return fmt.Errorf("invalid %s: %v", field.Tag.Get("yaml"), err)
			}
		}
	}

	// Checking frequency range logic
	min, _ := ParseFreq(c.MinFreq)
	max, _ := ParseFreq(c.MaxFreq)
	minAndmaxAreSet := (min != -1 && max != -1) && (min != 0 && max != 0)
	if (max < min) && minAndmaxAreSet {
		return fmt.Errorf(
			"max frequency (%d) cannot be less than min frequency (%d)",
			max, min)
	}
	if (min == max) && minAndmaxAreSet {
		return fmt.Errorf("min and max frequency cannot be the same: %d", min)
	}

	return nil
}

func ParseFreq(freq string) (int, error) {
	if freq == "" {
		return -1, nil // No frequency limits set, nothing to parse
	}
	if err := CheckFreqFormat(freq); err != nil {
		return -1, err
	}

	s := strings.ToLower(strings.TrimSpace(freq))
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
	}

	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return -1, fmt.Errorf("failed to parse frequency value: %v", err)
	}

	kHz := int(val * multiplier)
	return kHz, nil
}
