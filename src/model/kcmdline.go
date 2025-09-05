package model

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/canonical/rt-conf/src/cpulists"
)

var isolcpuFlags = []string{"domain", "nohz", "managed_irq"}

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

		if len(keyValue) == 0 || keyValue[0] == "" {
			return fmt.Errorf("empty parameter detected")
		}

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
				return fmt.Errorf("%q has an invalid value: %q: %v", key, value, err)
			}
		case "nohz":
			if value != "on" && value != "off" {
				return fmt.Errorf("%q must be set to either 'on' or 'off', got %q", key, value)
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
