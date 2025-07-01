package model

import (
	"fmt"
	"reflect"
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
				return fmt.Errorf("invalid cpus: %v", err)
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
