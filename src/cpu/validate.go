package cpu

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

var isolcpuFlags []string = []string{"domain", "nohz", "managed_irq"}

func ValidateList(s string) error {
	max, err := TotalAvailable()
	if err != nil {
		return fmt.Errorf("failed to get total available CPUs: %v", err)
	}
	return validateList(s, max)
}

func validateList(s string, max int) error {
	_, err := parseCPUs(s, max)
	return err
}

func validateListWithFlags(s string, f []string, max int) error {
	hasFlag := true
	// Split the string into two parts by the first comma
	parts := strings.SplitN(s, ",", 2)

	// If it converts to a number, it's not a flag
	_, err := strconv.Atoi(parts[0])

	// Check if the first part isn't a flag
	if len(parts) != 2 ||
		err == nil ||
		strings.Contains(parts[0], "-") ||
		strings.Contains(parts[0], ":") ||
		strings.Contains(parts[0], "/") {
		hasFlag = false
	}

	// If it's a flag, check if it's a valid flag
	if hasFlag {
		if !slices.Contains(f, parts[0]) {
			return fmt.Errorf("invalid flag: %s, expected one of %v",
				parts[0], f)
		}

		_, err := parseCPUs(parts[1], max)
		return err
	}

	_, err = parseCPUs(s, max)
	return err
}

func ValidateIsolCPUs(s string) error {
	max, err := TotalAvailable()
	if err != nil {
		return fmt.Errorf("failed to get total available CPUs: %v", err)
	}
	return validateListWithFlags(s, isolcpuFlags, max)
}
