package data

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/canonical/rt-conf/src/common"
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

// TODO: Validate mutual exclusive cpu lists

func (c IRQFilter) validateIRQField(name string, value string, tag string) error {
	fmt.Println("[DEBUG] Validating IRQ field")
	switch {
	case tag == "cpulist":
		num, err := GetHigherIRQ()
		fmt.Printf("[DEBUG] Higher IRQ num: %v\n", num)
		if err != nil {
			return err
		}
		_, err = cpu.ValidateCPUListSyntax(value, num)
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

func GetHigherIRQ() (int, error) {
	files, err := os.ReadDir(common.SysKernelIRQ)
	if err != nil {
		return 0, err
	}
	var irqs []int
	for _, file := range files {
		num, err := strconv.Atoi(file.Name())
		if err != nil {
			continue
		}
		if file.IsDir() {
			irqs = append(irqs, num)
		}
	}
	if len(irqs) == 0 {
		return 0, fmt.Errorf("no IRQs found")
	}
	bigger := irqs[0]
	for _, irq := range irqs {
		if irq > bigger {
			bigger = irq
		}
	}
	return bigger, nil
}
