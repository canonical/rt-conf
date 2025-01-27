package data

import (
	"fmt"
	"regexp"

	"github.com/canonical/rt-conf/src/cpu"
	"github.com/canonical/rt-conf/src/helpers"
)

type IRQTunning struct {
	CPUs   string    `yaml:"cpus"`
	Filter IRQFilter `yaml:"filter"`
}

func (c IRQTunning) Validate() error {
	err := c.Filter.Validate()
	if err != nil {
		return fmt.Errorf("IRQFilter validation failed: %v", err)
	}
	err = cpu.ValidateList(c.CPUs)
	if err != nil {
		return fmt.Errorf("invalid cpus: %v", err)
	}
	return nil
}

type IRQFilter struct {
	Number   string `yaml:"number" validation:"cpulist"`
	Action   string `yaml:"action" validation:"regex"`
	ChipName string `yaml:"chip_name" validation:"regex"`
	Name     string `yaml:"name" validation:"regex"`
	Type     string `yaml:"type" validation:"regex"`
}

type IRQs struct {
	IsolateCPU string `yaml:"remove-from-cpus"`
	IRQHandler string `yaml:"handle-on-cpus"`
}

func (c IRQFilter) Validate() error {
	return helpers.Validate(c, c.validateIRQField)
}

func (c IRQFilter) validateIRQField(name string, value string, tag string) error {
	switch {
	case tag == "cpulist":
		err := cpu.ValidateIsolCPUs(value)
		if err != nil {
			return fmt.Errorf("on field %v: invalid irq list: %v", name,
				err)
		}
	case tag == "regex":
		_, err := regexp.Compile(value)
		if err != nil {
			return fmt.Errorf("on field %v: invalid regex: %v", name, err)
		}
	default:
		return fmt.Errorf("on field %v: invalid tag: %v", name, tag)
	}
	return nil
}
