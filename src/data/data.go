package data

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"

	"github.com/canonical/rt-conf/src/cpu"
)

// TODO: this should be dropped in favor or parsing the default grub file
//
//	** NOTE: instead of using a hardcoded pattern
//	** we should source the /etc/default/grub file, since it's basically
//	** an environment file (key=value) and we can parse it with the
//	** `os` package
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
	return c.KernelCmdline.Validate()
}

type IRQs struct {
	IsolateCPU string `yaml:"remove-from-cpus"`
	IRQHandler string `yaml:"handle-on-cpus"`
}

// KernelCmdline represents the kernel command line options.
type KernelCmdline struct {
	// Isolate CPUs
	IsolCPUs string `yaml:"isolcpus" kcmd:"isolcpus" validation:"cpulist:isolcpus"`

	// Enable/Disable dyntick idle
	Nohz string `yaml:"nohz" kcmd:"nohz" validation:"oneof:on,off"`

	// CPUs for adaptive ticks
	NohzFull string `yaml:"nohz_full" kcmd:"nohz_full" "validation:"cpulist"`

	// CPUs for kthreads
	KthreadCPUs string `yaml:"kthread_cpus" kcmd:"kthread_cpus" validation:"cpulist"`

	// CPUs for IRQs
	IRQaffinity string `yaml:"irqaffinity" kcmd:"irqaffinity" validation:"cpulist"`
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

// Custom unmarshal function with validation
func (c KernelCmdline) Validate() error {
	// Validate fields based on struct tags
	v := reflect.ValueOf(c)
	t := reflect.TypeOf(c)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		tag := field.Tag.Get("validation")
		if tag == "" {
			continue // No validation tag, skip
		}

		// * * NOTE: For now it's okay to cast to string
		// * * since ther is only strings on KernelCmdline struct
		value, ok := v.Field(i).Interface().(string)
		if !ok {
			return fmt.Errorf("value for field %s is not a string", value)
		}

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
	case strings.HasPrefix(tag, "cpulist"):
		if strings.HasSuffix(tag, "isolcpus") {
			err := cpu.ValidateIsolCPUs(value)
			if err != nil {
				return fmt.Errorf("on field %v: invalid isolcpus: %v", name, err)
			}
		} else {
			err := cpu.ValidateList(value)
			if err != nil {
				return fmt.Errorf("on field %v: invalid cpulist: %v", name, err)
			}
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
