package cpulists

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseCPUs parses a CPU list into a set of integers, supporting all formats
// Inspired in the Kernel documentation:
// https://docs.kernel.org/admin-guide/kernel-parameters.html#cpu-lists
// I had nightmares about this function
func ParseCPUs(cpuList string) (CPUs, error) {
	max, err := TotalAvailable()
	if err != nil {
		return nil, fmt.Errorf("failed to get total available CPUs: %v", err)
	}
	return parseCPUs(cpuList, max)
}

func parseCPUs(cpuList string, totalCPUs int) (CPUs, error) {
	cpus := make(CPUs)
	items := strings.Split(cpuList, ",")

	for _, item := range items {
		item = strings.TrimSpace(item)

		// Handle "all"
		if item == "all" {
			for i := 0; i < totalCPUs; i++ {
				cpus[i] = true
			}
			continue
		}

		// Handle "N" or "n"
		item = strings.ReplaceAll(item, "N", strconv.Itoa(totalCPUs-1))
		item = strings.ReplaceAll(item, "n", strconv.Itoa(totalCPUs-1))

		if strings.Contains(item, ":") {
			if err := handleCPUGroup(item, cpus, totalCPUs); err != nil {
				return nil, err
			}
			continue
		}

		if strings.Contains(item, "-") {
			if err := handleCPURange(item, cpus, totalCPUs); err != nil {
				return nil, err
			}
			continue
		}

		if err := handleSingleCPU(item, cpus, totalCPUs); err != nil {
			return nil, err
		}
	}
	return cpus, nil
}

// Handle the format:
// <cpu number>-<cpu number>:<used size>/<group size>
func handleCPUGroup(item string, cpus CPUs, t int) error {
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
	if usedSize < 1 {
		return fmt.Errorf("used size must be at least 1, got: %s", groupParts[0])
	}
	if usedSize > t {
		return fmt.Errorf("used size greater than total CPUs: %s", groupParts[0])
	}
	groupSize, err := strconv.Atoi(groupParts[1])
	if err != nil {
		return fmt.Errorf("invalid group size: %s", groupParts[1])
	}

	// Split the range into groups and take the first "usedSize" CPUs
	for i := start; i <= end; i += groupSize {
		for j := 0; j < usedSize && (i+j) <= end; j++ {
			cpus[i+j] = true
		}
	}
	return nil
}

// Handle: <cpu number>-<cpu number>
func handleCPURange(item string, cpus CPUs, t int) error {
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
		return fmt.Errorf("invalid end of range: %s", parts[1])
	}
	if end >= t {
		return fmt.Errorf("end of range greater than total CPUs: %s", item)
	}
	if start > end {
		return fmt.Errorf("start of range greater than end: %s", item)
	}
	for i := start; i <= end; i++ {
		cpus[i] = true
	}
	return nil
}

// Handle <cpu number>
func handleSingleCPU(item string, cpus CPUs, t int) error {
	// Handle single CPU
	cpu, err := strconv.Atoi(item)
	if err != nil {
		return fmt.Errorf("invalid CPU: %s", item)
	}
	if cpu >= t {
		return fmt.Errorf("CPU greater than total CPUs: %s", item)
	}
	cpus[cpu] = true
	return nil
}
