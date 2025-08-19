package model

import (
	"fmt"
)

type InternalConfig struct {
	Data Config

	GrubCfg Grub
}

type (
	PwrMgmt    map[string]CpuGovernanceRule
	Interrupts map[string]IRQTuning
)

type Config struct {
	KernelCmdline KernelCmdline `yaml:"kernel-cmdline"`
	Interrupts    Interrupts    `yaml:"irq-tuning"`
	CpuGovernance PwrMgmt       `yaml:"cpu-governance"`
}

func (c Config) Validate() error {
	err := c.KernelCmdline.Validate()
	if err != nil {
		return fmt.Errorf("failed to validate kernel cmdline: %v", err)
	}
	for _, irq := range c.Interrupts {
		err := irq.Validate()
		if err != nil {
			return fmt.Errorf("failed to validate irq tuning: %v", err)
		}
	}

	for label, pwrprof := range c.CpuGovernance {
		err := pwrprof.Validate()
		if err != nil {
			return fmt.Errorf(
				"failed to validate cpu governance rule #%s: %s", label, err)
		}
	}

	return nil
}
