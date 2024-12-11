package interrupts

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseCPUs parses a CPU list into a set of integers, supporting all formats
// Inspired in the Kernel documentation:
// https://docs.kernel.org/admin-guide/kernel-parameters.html#cpu-lists
// I had nightmares about this function
func ParseCPUs(cpuList string, totalCPUs int) (Cores, error) {
	cpus := make(Cores)
	items := strings.Split(cpuList, ",")

	for _, item := range items {
		item = strings.TrimSpace(item)

		// Handle "all"
		if item == "all" {
			for i := 0; i < totalCPUs; i++ {
				cpus[i] = struct{}{}
			}
			continue
		}

		// Handle "N" or "n"
		item = strings.ReplaceAll(item, "N", strconv.Itoa(totalCPUs-1))
		item = strings.ReplaceAll(item, "n", strconv.Itoa(totalCPUs-1))

		if strings.Contains(item, ":") {
			if err := HandleFormatedLists(item, cpus, totalCPUs); err != nil {
				return nil, err
			}
			continue
		}

		if strings.Contains(item, "-") {
			if err := HandleCPURange(item, cpus, totalCPUs); err != nil {
				return nil, err
			}
			continue
		}

		if err := HandleSingleCPU(item, cpus, totalCPUs); err != nil {
			return nil, err
		}
	}
	return cpus, nil
}

// Handle the format:
// <cpu number>-<cpu number>:<used size>/<group size>
func HandleFormatedLists(item string, cpus Cores, t int) error {
	rangePart, groupPart, found := strings.Cut(item, ":")
	if !found {
		return fmt.Errorf("invalid format: %s", item)
	}
	startEnd := strings.Split(rangePart, "-")
	if len(startEnd) != 2 {
		return fmt.Errorf("invalid range: %s", rangePart)
	}
	start, err := strconv.Atoi(startEnd[0])
	if err != nil {
		return fmt.Errorf("invalid start of range: %s", startEnd[0])
	}
	end, err := strconv.Atoi(startEnd[1])
	if err != nil {
		return fmt.Errorf("invalid end of range: %s", startEnd[1])
	}
	if end >= t {
		return fmt.Errorf("end of range greater than total CPUs: %s", item)
	}

	groupParts := strings.Split(groupPart, "/")
	if len(groupParts) != 2 {
		return fmt.Errorf("invalid group size or used size: %s", groupPart)
	}
	usedSize, err := strconv.Atoi(groupParts[0])
	if err != nil {
		return fmt.Errorf("invalid used size: %s", groupParts[0])
	}
	groupSize, err := strconv.Atoi(groupParts[1])
	if err != nil {
		return fmt.Errorf("invalid group size: %s", groupParts[1])
	}

	// Split the range into groups and take the first "usedSize" CPUs
	for i := start; i <= end; i += groupSize {
		for j := 0; j < usedSize && (i+j) <= end; j++ {
			cpus[i+j] = struct{}{}
		}
	}
	return nil
}

// Handle: <cpu number>-<cpu number>
func HandleCPURange(item string, cpus Cores, t int) error {
	parts := strings.Split(item, "-")
	if len(parts) != 2 {
		return fmt.Errorf("invalid range: %s", item)
	}
	start, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("invalid start of range: %s", parts[0])
	}
	end, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid end of range: %s", parts[2])
	}
	if end >= t {
		return fmt.Errorf("end of range greater than total CPUs: %s", item)
	}
	if start > end {
		return fmt.Errorf("start of range greater than end: %s", item)
	}
	for i := start; i <= end; i++ {
		cpus[i] = struct{}{}
	}
	return nil
}

// Handle <cpu number>
func HandleSingleCPU(item string, cpus Cores, t int) error {
	// Handle single CPU
	cpu, err := strconv.Atoi(item)
	if err != nil {
		return fmt.Errorf("invalid CPU: %s", item)
	}
	if cpu >= t {
		return fmt.Errorf("CPU greater than total CPUs: %s", item)
	}
	cpus[cpu] = struct{}{}
	return nil
}

// CheckMutExclusive:
// checks if two CPU lists are mutually exclusive
func CheckMutExclusive(list1, list2 string, totalCPUs int) (bool, error) {
	set1, err := ParseCPUs(list1, totalCPUs)
	if err != nil {
		return false, err
	}
	set2, err := ParseCPUs(list2, totalCPUs)
	if err != nil {
		return false, err
	}

	for cpu := range set1 {
		if _, exists := set2[cpu]; exists {
			return false, nil
		}
	}

	return true, nil
}
