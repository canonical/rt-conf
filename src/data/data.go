package data

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"

	"github.com/canonical/rt-conf/src/cpu"
	"gopkg.in/yaml.v3"
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

func (c *Config) UnmarshalYAML(node *yaml.Node) error {

	type rawConfig Config
	var raw rawConfig
	// Decode YAML normally into KernelCmdline
	if err := node.Decode(&raw); err != nil {
		return err
	}
	*c = Config(raw)

	return nil
}

type IRQs struct {
	IsolateCPU string `yaml:"remove-from-cpus"`
	IRQHandler string `yaml:"handle-on-cpus"`
}

// KernelCmdline represents the kernel command line options.
type KernelCmdline struct {
	// Isolate CPUs
	IsolCPUs string `yaml:"isolcpus" validation:"cpulist"`

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
func (k *KernelCmdline) UnmarshalYAML(node *yaml.Node) error {
	// Decode YAML normally
	type rawKernelCmdline KernelCmdline // To avoid recursion
	var raw rawKernelCmdline
	if err := node.Decode(&raw); err != nil {
		return err
	}

	// Validate fields based on struct tags
	v := reflect.ValueOf(raw)
	t := reflect.TypeOf(raw)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		tag := field.Tag.Get("validation")
		if tag == "" {
			continue // No validation tag, skip
		}

		// * * NOTE: For now it's okay to cast to string
		// * * since ther is only strings on KernelCmdline struct
		value := v.Field(i).Interface().(string)

		if value == "" {
			continue
		}

		err := validateField(field.Name, value, tag)
		if err != nil {
			return fmt.Errorf("validation failed for field %s: %v",
				field.Name, err)
		}
	}

	// Assign validated values
	*k = KernelCmdline(raw)
	return nil
}

func (k *Config) Marshal(in interface{}) (out []byte, err error) {
	return yaml.Marshal(in)
}

// validateField validates a field based on the tag
func validateField(name string, value interface{}, tag string) error {
	parts := strings.Split(tag, ",")
	for _, rule := range parts {
		switch {
		case rule == "cpulist":
			err := cpu.ValidateList(value.(string))
			if err != nil {
				return fmt.Errorf("on field %v: invalid cpulist: %v", name, err)
			}
		case strings.HasPrefix(rule, "oneof:"):
			options := strings.Split(rule[len("oneof:"):], "|")
			if !slices.Contains(options, fmt.Sprintf("%v", value)) {
				return fmt.Errorf("value must be one of %v", options)
			}
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
