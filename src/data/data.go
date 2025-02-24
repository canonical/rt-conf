package data

import (
	"fmt"
)

type InternalConfig struct {
	CfgFile string
	Data    Config

	GrubDefault Grub
}

type Config struct {
	Interrupts    []IRQTunning        `yaml:"irq_tunning"`
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
			return fmt.Errorf("failed to validate irq tunning: %v", err)
		}
	}

	for i, pwrprof := range c.CpuGovernance {
		err := pwrprof.Validate()
		if err != nil {
			return fmt.Errorf(
				"failed to validate cpu governance rule #%v: %v", (i + 1), err)
		}
	}

	return nil
}
