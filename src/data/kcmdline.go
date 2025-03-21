package data

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/canonical/rt-conf/src/cpulists"
	"github.com/canonical/rt-conf/src/helpers"
)

var isolcpuFlags = []string{"domain", "nohz", "managed_irq"}

type Params map[string]string

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
			_, _, err := cpulists.ParseWithFlags(value, isolcpuFlags)
			if err != nil {
				return fmt.Errorf("on field %v: invalid isolcpus: %v", name,
					err)
			}
		} else {
			_, err := cpulists.Parse(value)
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

func ConstructKeyValuePairs(v *KernelCmdline) (Params, error) {
	kvpairs := make(Params, 0)
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
		kvpairs[key] = value
	}
	return kvpairs, nil
}

func CmdlineToParams(cmdline string) Params {
	kvpairs := make(Params)
	for _, p := range strings.Split(cmdline, " ") {
		pair := strings.Split(p, "=")
		// Value is optional for some kernel cmdline parameters
		if len(pair) != 2 {
			kvpairs[p] = ""
			continue
		}
		kvpairs[pair[0]] = pair[1]
	}
	return kvpairs
}

func ParamsToCmdline(params Params) string {
	var kcmds []string
	for k, v := range params {
		if v != "" && k != "" {
			kcmds = append(kcmds, fmt.Sprintf("%s=%s", k, v))
		}
		// Handle the case for parameters without a value such
		// as "quiet" and "splash"
		if v == "" && k != "" {
			kcmds = append(kcmds, k)
		}
	}
	return strings.Join(kcmds, " ")
}
