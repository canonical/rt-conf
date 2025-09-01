package model

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/canonical/rt-conf/src/cpulists"
)

var isolcpuFlags = []string{"domain", "nohz", "managed_irq"}

// Params represents key-value pairs for kernel parameters.
type Params map[string]string

// KernelCmdline represents the kernel command line options.
type KernelCmdline struct {
	Parameters []string `yaml:"parameters"`
}

const (
	// Maximum kernel command line characters length per architecture.
	// See COMMAND_LINE_SIZE macro in kernel source code:
	// - arch/x86/include/asm/setup.h
	// - arch/arm64/include/uapi/asm/setup.h
	// The value 2048 is defined for amd64 and arm64
	CommandLineSize = 2048
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
	return KernelCmdline{Parameters: fields}
}

// ToParams converts the KernelCmdline into Params map.
// If a parameter has no explicit value, the value is set to "".
func (k KernelCmdline) ToParams() Params {
	params := make(Params, len(k.Parameters))

	// Parse each parameter

	for _, param := range k.Parameters {
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

// validateParameterFormat performs syntax validation by checking kernel parameter formatting rules
func (k KernelCmdline) validateParameterFormat() error {
	totalLen := 0
	for _, p := range k.Parameters {
		// Total length includes spaces between parameters so +1 for each param
		// unless it's the last one, but we are not checking that here, which gives us
		// a pratical limit of CommandLineSize -1 characters.
		totalLen += len(p) + 1
		if totalLen > CommandLineSize {
			return fmt.Errorf("command line exceeds maximum length of %d bytes", CommandLineSize)
		}

		keyValue := strings.SplitN(p, "=", 2)
		key := keyValue[0]

		if !validName.MatchString(key) {
			return fmt.Errorf("invalid parameter name: %q", key)
		}
	}

	return nil
}

// Validate performs comprehensive validation
func (k KernelCmdline) Validate() error {
	if err := k.validateParameterFormat(); err != nil {
		return err
	}
	return k.validateParameterValues()
}

// validateParameterValues performs semantic validation on known parameters for specific rules
func (k KernelCmdline) validateParameterValues() error {
	for _, p := range k.Parameters {
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

		// Validate parameters based on key
		switch key {
		case "kthread_cpus", "irqaffinity", "rcu_nocbs", "nohz_full":
			if _, err := cpulists.Parse(value); err != nil {
				return fmt.Errorf("%q does not contain a valid CPU List: %q: %v", key, value, err)
			}
		case "isolcpus":
			if _, _, err := cpulists.ParseWithFlags(value, isolcpuFlags); err != nil {
				return fmt.Errorf("parameter %q has invalid isolcpus %q: %v", key, value, err)
			}
		case "nohz":
			if value != "on" && value != "off" {
				return fmt.Errorf("parameter %q value must be 'on' or 'off', got %q", key, value)
			}
		default:
			log.Printf("Warning: Parameter %q not recognized by rt-conf; skipping specific validation", key)
		}
	}
	return nil
}

// HasDuplicates checks for duplicate parameters with different values
func (k KernelCmdline) HasDuplicates() error {
	params := make(map[string]string)

	for _, p := range k.Parameters {
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

// MergeWithOrderPreservation merges existing and new parameters while preserving order.
// Existing parameters come first in their original order, then new parameters in the order they appear in newCmdline.
// If a parameter exists in both, the new value overwrites the existing one but keeps the original position.
func MergeWithOrderPreservation(existingCmdline, newCmdline KernelCmdline) string {
	if len(existingCmdline.Parameters) == 0 && len(newCmdline.Parameters) == 0 {
		return ""
	}

	newParams := newCmdline.ToParams()

	// Track which new parameters we've already processed
	processedNewParams := make(map[string]bool)

	var result []string

	// First, process existing parameters
	for _, param := range existingCmdline.Parameters {
		if param == "" {
			continue
		}

		var key string
		if idx := strings.Index(param, "="); idx != -1 {
			key = param[:idx]
		} else {
			key = param
		}

		// If this parameter is being overridden by new config, use the new value
		if newValue, exists := newParams[key]; exists {
			if newValue != "" {
				result = append(result, fmt.Sprintf("%s=%s", key, newValue))
			} else {
				result = append(result, key)
			}
			processedNewParams[key] = true
		} else {
			// Keep the existing parameter as-is
			result = append(result, param)
		}
	}

	// Then, append any new parameters that weren't already processed (in their original order)
	for _, param := range newCmdline.Parameters {
		if param == "" {
			continue
		}

		var key string
		if idx := strings.Index(param, "="); idx != -1 {
			key = param[:idx]
		} else {
			key = param
		}

		if !processedNewParams[key] {
			result = append(result, param)
		}
	}

	return strings.Join(result, " ")
}

// Merge merges other Params into this one, overwriting existing keys.
func (p Params) Merge(other Params) {
	for k, v := range other {
		p[k] = v
	}
}
