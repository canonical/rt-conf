// This package provides a parser for CPU Lists strings.
// It supports all formats described in the Kernel documentation:
// https://docs.kernel.org/admin-guide/kernel-parameters.html#cpu-lists

package cpulists

import (
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"
)

type CPUs map[int]bool

// Parse parses a CPU Lists string into CPUs map
// It performs parsing based on the total number of available CPUs
func Parse(cpuLists string) (CPUs, error) {
	total, err := totalCPUs()
	if err != nil {
		return nil, fmt.Errorf("failed to get total available CPUs: %v", err)
	}
	return ParseForCPUs(cpuLists, total)
}

// ParseForCPUs parses a CPU Lists string into CPUs map
func ParseForCPUs(cpuLists string, totalCPUs int) (CPUs, error) {
	cpus := make(CPUs)
	items := strings.Split(cpuLists, ",")

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

// ParseWithFlagsForCPUs parses a CPU Lists string into CPUs map
func ParseWithFlagsForCPUs(cpuLists string, validFlags []string, totalCPUs int) (CPUs, string, error) {
	hasFlag := true
	// Split the string into two parts by the first comma
	parts := strings.SplitN(cpuLists, ",", 2)

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
		if !slices.Contains(validFlags, parts[0]) {
			return nil, "", fmt.Errorf("invalid flag: %s, expected one of %v", parts[0], validFlags)
		}

		cpus, err := ParseForCPUs(parts[1], totalCPUs)
		if err != nil {
			return nil, "", fmt.Errorf("failed to parse CPUs: %s", err)
		}
		return cpus, parts[0], nil
	}

	cpus, err := ParseForCPUs(cpuLists, totalCPUs)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse CPUs: %s", err)
	}
	return cpus, "", nil
}

// ParseWithFlagsForCPUs parses a CPU Lists string into CPUs map
// It performs parsing based on the total number of available CPUs
func ParseWithFlags(cpuLists string, validFlags []string) (CPUs, string, error) {
	totalCPUs, err := totalCPUs()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get total available CPUs: %v", err)
	}

	return ParseWithFlagsForCPUs(cpuLists, validFlags, totalCPUs)
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

func GenCPUlist(cpus []int) string {
	if len(cpus) == 0 {
		return ""
	}
	list := deduplicateCPUs(cpus)

	var parts []string
	start := list[0]
	end := start

	for i := 1; i < len(list); i++ {
		if list[i] == end+1 {
			end = list[i]
		} else {
			if start == end {
				parts = append(parts, fmt.Sprintf("%d", start))
			} else {
				parts = append(parts, fmt.Sprintf("%d-%d", start, end))
			}
			start = list[i]
			end = start
		}
	}
	// Final range
	if start == end {
		parts = append(parts, fmt.Sprintf("%d", start))
	} else {
		parts = append(parts, fmt.Sprintf("%d-%d", start, end))
	}

	return strings.Join(parts, ",")
}

func deduplicateCPUs(cpus []int) (cpulist []int) {
	cpuMap := make(CPUs)
	for _, cpu := range cpus {
		cpuMap[cpu] = true
	}
	for cpu := range cpuMap {
		cpulist = append(cpulist, cpu)
	}
	sort.Ints(cpulist)
	return cpulist
}
