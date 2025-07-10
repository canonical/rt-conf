package model

import (
	"fmt"
)

type InternalConfig struct {
	Data Config

	GrubCfg Grub
}

type PwrMgmt map[string]CpuGovernanceRule
type Interrupts map[string]IRQTuning

type Config struct {
	Interrupts    Interrupts    `yaml:"irq_tuning"`
	KernelCmdline KernelCmdline `yaml:"kernel_cmdline"`
	CpuGovernance PwrMgmt       `yaml:"cpu_governance"`
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
