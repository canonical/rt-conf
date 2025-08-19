package model

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/canonical/rt-conf/src/cpulists"
)

var isolcpuFlags = []string{"domain", "nohz", "managed_irq"}

// Params represents key-value pairs for kernel parameters.
type Params map[string]string

// KernelCmdline represents the kernel command line options.
type KernelCmdline []string

const (
	// Maximum kernel command line length per architecture.
	// See COMMAND_LINE_SIZE macro, 2048 for amd64 and arm64
	// See in kernel source code:
	// - arch/x86/include/asm/setup.h
	// - arch/arm64/include/uapi/asm/setup.h
	COMMAND_LINE_SIZE = 2048
)

// Regex for valid parameter names
// - must start with letters
// - can contain: letters, digits, underscores, dots, hyphens
var validName = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9\._\-]*`)

// NewKernelCmdline creates a new KernelCmdline from a string.
func NewKernelCmdline(cmdline string) KernelCmdline {
	if cmdline == "" {
		return KernelCmdline{}
	}
	fields := strings.Fields(cmdline)
	return KernelCmdline(fields)
}

// Join returns the kernel command line as a single string.
func (k KernelCmdline) Join() string {
	return strings.Join(k, " ")
}

// ToParams converts the KernelCmdline into key=value pairs.
// If a parameter has no explicit value, the value is set to "".
func (k KernelCmdline) ToParams() Params {
	params := make(Params, len(k))

	for _, param := range k {
		if param == "" {
			continue
		}

		// Split on first '=' only to handle values that contain '='
		if idx := strings.Index(param, "="); idx != -1 {
			key := param[:idx]
			value := param[idx+1:]
			params[key] = value
		} else {
			// Parameter without value (e.g., "quiet", "splash")
			params[param] = ""
		}
	}

	return params
}

// ValidateKernelParams checks kernel parameter formatting rules
func (k KernelCmdline) ValidateKernelParams() error {
	totalLen := 0
	for i, p := range k {
		totalLen += len(p) + 1
		if totalLen > COMMAND_LINE_SIZE {
			return fmt.Errorf("command line exceeds maximum length of %d bytes", COMMAND_LINE_SIZE)
		}

		keyValue := strings.SplitN(p, "=", 2)
		key := keyValue[0]

		if !validName.MatchString(key) {
			return fmt.Errorf("invalid parameter name at index %d: %q", i, key)
		}
	}

	return nil
}

// Validate performs comprehensive validation
func (k KernelCmdline) Validate() error {
	if err := k.ValidateKernelParams(); err != nil {
		return err
	}
	return k.validateKnownParams()
}

func (k KernelCmdline) validateKnownParams() error {
	cpulistParams := []string{
		"kthread_cpus",
		"irqaffinity",
		"rcu_nocbs",
	}

	for _, p := range k {
		// Skip empty parameters
		if p == "" {
			continue
		}

		// Split parameter into key and value
		parts := strings.SplitN(p, "=", 2)
		if len(parts) < 2 {
			// No value, skip validation for boolean parameters
			continue
		}

		key := parts[0]
		value := parts[1]

		// Validate cpulist parameters
		for _, cpuParam := range cpulistParams {
			if key == cpuParam {
				if _, err := cpulists.Parse(value); err != nil {
					return fmt.Errorf("parameter %q has invalid cpulist %q: %v", key, value, err)
				}
			}
		}
		// Handle isolcpus parameter
		if key == "isolcpus" {
			if _, _, err := cpulists.ParseWithFlags(value, isolcpuFlags); err != nil {
				return fmt.Errorf("parameter %q has invalid isolcpus %q: %v", key, value, err)
			}
		}
		// Handle nohz parameter
		if key == "nohz" {
			if value != "on" && value != "off" {
				return fmt.Errorf("parameter %q value must be 'on' or 'off', got %q", key, value)
			}
		}
	}
	return nil
}

// HasDuplicates checks for duplicate parameters with different values
func (k KernelCmdline) HasDuplicates() error {
	params := make(map[string]string)

	for _, p := range k {
		if p == "" {
			continue
		}

		var key, value string
		if idx := strings.Index(p, "="); idx != -1 {
			key = p[:idx]
			value = p[idx+1:]
		} else {
			key = p
			value = ""
		}

		if existingValue, exists := params[key]; exists {
			// Allow duplicate parameters with the same value
			if existingValue != value {
				return fmt.Errorf("duplicate parameter %q with different values: %q and %q", key, existingValue, value)
			}
		} else {
			params[key] = value
		}
	}
	return nil
}

// CmdlineToParamsError represents an error in command line parsing.
type CmdlineToParamsError struct {
	Type string
	Msg  string
}

func (e *CmdlineToParamsError) Error() string {
	return fmt.Sprintf("cmdline parsing error: %s (type: %s)", e.Msg, e.Type)
}

// CmdlineToParams converts various command line representations into Params.
// Supported types: string, []string, KernelCmdline
func CmdlineToParams(cmdline any) (Params, error) {
	switch v := cmdline.(type) {
	case string:
		return NewKernelCmdline(v).ToParams(), nil
	case []string:
		return KernelCmdline(v).ToParams(), nil
	case KernelCmdline:
		return v.ToParams(), nil
	default:
		return nil, &CmdlineToParamsError{
			Type: reflect.TypeOf(cmdline).String(),
			Msg:  "unsupported type",
		}
	}
}

// ParamsToCmdline converts Params back to a command line string.
// Parameters are sorted alphabetically for consistent output.
func ParamsToCmdline(params Params) string {
	if len(params) == 0 {
		return ""
	}

	// Sort keys for consistent output
	keys := make([]string, 0, len(params))
	for k := range params {
		if k != "" { // Skip empty keys
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var parts []string
	for _, key := range keys {
		value := params[key]
		if value != "" {
			parts = append(parts, fmt.Sprintf("%s=%s", key, value))
		} else {
			// Parameter without value (e.g., "quiet", "splash")
			parts = append(parts, key)
		}
	}

	return strings.Join(parts, " ")
}

// Merge merges other Params into this one, overwriting existing keys.
func (p Params) Merge(other Params) {
	for k, v := range other {
		p[k] = v
	}
}
