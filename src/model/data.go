package model

import (
	"fmt"
)

type InternalConfig struct {
	CfgFile string
	Data    Config

	GrubDefault Grub
}

type Config struct {
	Interrupts    []IRQTuning         `yaml:"irq_tuning"`
	KernelCmdline KernelCmdline       `yaml:"kernel_cmdline"`
	CpuGovernance []CpuGovernanceRule `yaml:"cpu_governance"`
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

	for i, pwrprof := range c.CpuGovernance {
		err := pwrprof.Validate()
		if err != nil {
			return fmt.Errorf(
				"failed to validate cpu governance rule #%d: %s", (i + 1), err)
		}
	}

	return nil
}
