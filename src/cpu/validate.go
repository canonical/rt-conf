package cpu

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func ValidateList(s string) error {

	max, err := TotalAvailable()
	if err != nil {
		return fmt.Errorf("failed to get total available CPUs: %v", err)
	}
	_, err = ParseCPUs(s, max)
	return err
}

func ValidateListWithFlags(s string, f []string) error {

	max, err := TotalAvailable()
	if err != nil {
		return fmt.Errorf("failed to get total available CPUs: %v", err)
	}

	return validateListWithFlags(s, f, max)
}

func validateListWithFlags(s string, f []string, max int) error {

	hasFlag := true

	// Split the string into two parts by the first comma
	parts := strings.SplitN(s, ",", 2)

	_, err := strconv.Atoi(parts[0])

	// Check if the first part isn't a flag
	if len(parts) != 2 ||
		err == nil ||
		strings.Contains(parts[0], "-") ||
		strings.Contains(parts[0], ":") ||
		strings.Contains(parts[0], "/") {
		hasFlag = false
	}

	if hasFlag && err != nil && !slices.Contains(f, parts[0]) {
		return fmt.Errorf("invalid flag: %s", parts[0])
	}

	var errCPU error
	if hasFlag {
		_, errCPU = ParseCPUs(parts[1], max)
	} else if len(parts) == 2 {
		_, errCPU = ParseCPUs(parts[0], max)
	}

	return errCPU
}
