package model

import (
	"fmt"
	"regexp"
)

type InternalConfig struct {
	Data Config

	GrubCfg Grub
}

type PwrMgmt map[string]CpuGovernanceRule
type Interrupts map[string]IRQTuning

type Config struct {
	KernelCmdline KernelCmdline `yaml:"kernel-cmdline"`
	Interrupts    Interrupts    `yaml:"irq-tuning"`
	CpuGovernance PwrMgmt       `yaml:"cpu-governance"`
}

// Regex for valid snap options from snapd:
// See: https://github.com/canonical/snapd/blob/2.71/overlord/configstate/config/helpers.go#L36
var validRuleName = regexp.MustCompile("^(?:[a-z0-9]+-?)*[a-z](?:-?[a-z0-9])*$")

func (c Config) Validate() error {
	err := c.KernelCmdline.Validate()
	if err != nil {
		return fmt.Errorf("failed to validate kernel cmdline: %v", err)
	}
	for label, irq := range c.Interrupts {
		if !validRuleName.MatchString(label) {
			return fmt.Errorf("invalid rule name: %q", label)
		}
		err := irq.Validate()
		if err != nil {
			return fmt.Errorf("failed to validate irq tuning: %v", err)
		}
	}

	for label, pwrprof := range c.CpuGovernance {
		if !validRuleName.MatchString(label) {
			return fmt.Errorf("invalid rule name: %q", label)
		}
		err := pwrprof.Validate()
		if err != nil {
			return fmt.Errorf(
				"failed to validate cpu governance rule #%s: %s", label, err)
		}
	}

	return nil
}
