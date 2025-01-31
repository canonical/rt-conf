package data

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/canonical/rt-conf/src/cpu"
	"github.com/canonical/rt-conf/src/helpers"
)

// KernelCmdline represents the kernel command line options.
type KernelCmdline struct {
	// Isolate CPUs
	IsolCPUs string `yaml:"isolcpus" kcmd:"isolcpus"  validation:"cpulist:isolcpus"`

	// Enable/Disable dyntick idle
	Nohz string `yaml:"nohz" kcmd:"nohz" validation:"oneof:on,off"`

	// CPUs for adaptive ticks
	NohzFull string `yaml:"nohz_full" kcmd:"nohz_full" validation:"cpulist"`

	// CPUs for kthreads
	KthreadCPUs string `yaml:"kthread_cpus" kcmd:"kthread_cpus" validation:"cpulist"`

	// CPUs for IRQs
	IRQaffinity string `yaml:"irqaffinity" kcmd:"irqaffinity" validation:"cpulist"`
}

// Custom unmarshal function with validation
func (c KernelCmdline) Validate() error {
	return helpers.Validate(c, c.fieldValidator)
}

func (c KernelCmdline) fieldValidator(name string,
	value string, tag string) error {

	switch {
	case strings.HasPrefix(tag, "cpulist"):
		if strings.HasSuffix(tag, "isolcpus") {
			err := cpu.ValidateIsolCPUs(value)
			if err != nil {
				return fmt.Errorf("on field %v: invalid isolcpus: %v", name,
					err)
			}
		} else {
			err := cpu.ValidateList(value)
			if err != nil {
				return fmt.Errorf("on field %v: invalid cpulist: %v", name,
					err)
			}
		}

	case strings.HasPrefix(tag, "oneof:"):
		options := strings.Split(tag[len("oneof:"):], ",")
		if !slices.Contains(options, value) {
			return fmt.Errorf("value must be one of %v", options)
		}

	default:
		return fmt.Errorf("unhandled tag: %v", tag)

	}
	return nil
}

func ConstructKeyValuePairs(v *KernelCmdline) ([]string, error) {
	var keyValuePairs []string

	val := reflect.TypeOf(v)
	valValue := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		valValue = valValue.Elem()
	}
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		key := field.Tag.Get("kcmd")
		value := valValue.Field(i).String()
		if key == "" || value == "" {
			continue
		}

		keyValuePairs = append(keyValuePairs, fmt.Sprintf("%s=%s", key, value))
	}
	return keyValuePairs, nil
}
