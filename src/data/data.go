package data

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"

	"github.com/canonical/rt-conf/src/cpu"
)

var PatternGrubDefault = regexp.MustCompile(`^(GRUB_CMDLINE_LINUX=")([^"]*)(")$`)

type InternalConfig struct {
	CfgFile string
	Data    Config

	GrubDefault Grub
}

// FileTransformer interface with a TransformLine method.
// This method is used to transform a line of a file.
//
// NOTE: This interface can be implemented also for RPi on classic
type FileTransformer interface {
	TransformLine(string) string
	GetFilePath() string
	GetPattern() *regexp.Regexp
}

type Core interface {
	InjectToFile(pattern *regexp.Regexp) error
}

type Grub struct {
	File    string
	Pattern *regexp.Regexp
}

type Config struct {
	Interrupts    IRQs          `yaml:"interrupts"`
	KernelCmdline KernelCmdline `yaml:"kernel_cmdline"`
}

func (c Config) Validate() error {
	fmt.Println("\n\n----------------[DEBUG]Validating config----------------")
	return c.KernelCmdline.Validate()
}

type IRQs struct {
	IsolateCPU string `yaml:"remove-from-cpus"`
	IRQHandler string `yaml:"handle-on-cpus"`
}

// KernelCmdline represents the kernel command line options.
type KernelCmdline struct {
	// Isolate CPUs
	IsolCPUs string `yaml:"isolcpus" validation:"isolcpus"`

	// Enable/Disable dyntick idle
	Nohz string `yaml:"nohz" validation:"oneof:on,off"`

	// CPUs for adaptive ticks
	NohzFull string `yaml:"nohz_full" validation:"cpulist"`

	// CPUs for kthreads
	KthreadCPUs string `yaml:"kthread_cpus" validation:"cpulist"`

	// CPUs for IRQs
	IRQaffinity string `yaml:"irqaffinity" validation:"cpulist"`
}

// Custom unmarshal function with validation
func (c KernelCmdline) Validate() error {
	// Validate fields based on struct tags
	v := reflect.ValueOf(c)
	t := reflect.TypeOf(c)

	fmt.Println("\n[DEBUG] Value of KernelCmdline: ", v)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		tag := field.Tag.Get("validation")
		fmt.Println("\n[DEBUG] Tag: ", tag)
		if tag == "" {
			continue // No validation tag, skip
		}

		// * * NOTE: For now it's okay to cast to string
		// * * since ther is only strings on KernelCmdline struct
		value, ok := v.Field(i).Interface().(string)
		if !ok {
			return fmt.Errorf("value for field %s is not a string", value)
		}
		fmt.Println("\n[DEBUG] Value: ", value)

		if value == "" {
			continue
		}
		err := validateField(field.Name, value, tag)
		if err != nil {
			return fmt.Errorf("validation failed for field %s: %v",
				field.Name, err)
		}
	}
	return nil
}

func validateField(name string, value string, tag string) error {
	switch {
	case tag == "cpulist":
		err := cpu.ValidateList(value)
		if err != nil {
			return fmt.Errorf("on field %v: invalid cpulist: %v", name, err)
		}

	case tag == "isolcpus":
		flags := []string{
			"domain",
			"nohz",
			"managed_irq",
		}
		err := cpu.ValidateListWithFlags(value, flags)
		if err != nil {
			return fmt.Errorf("on field %v: invalid isolcpus: %v", name, err)
		}

	case strings.HasPrefix(tag, "oneof:"):
		options := strings.Split(tag[len("oneof:"):], ",")
		if !slices.Contains(options, value) {
			return fmt.Errorf("value must be one of %v", options)
		}

	}
	return nil
}

// Param defines a mapping between a YAML field and a kernel parameter.
type Param struct {
	YAMLName    string
	CmdlineName string
	TransformFn func(interface{}) string // Function to transform value to string
}
