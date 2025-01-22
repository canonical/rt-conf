package cpu

import (
	"fmt"
	"strconv"
	"strings"
)

type CPUs map[int]bool

// Generate Complement (list') given the list
//
// returns list'
func GenerateComplementCPUList(list string, maxcpus int) (string, error) {
	var listprime []string
	cpus, err := ParseCPUs(list, maxcpus)
	if err != nil {
		return "", err
	}
	// OBS: This doesn't generate the most optimized CPU list
	// (only generates comma separated values) but it's good enough for now
	for i := 0; i < maxcpus; i++ {
		if _, exists := cpus[i]; !exists {
			listprime = append(listprime, strconv.Itoa(i))
		}
	}
	if len(listprime) == 0 {
		return "", fmt.Errorf(
			"unable to complement cpu list '%s' for total of %d CPUs", list, maxcpus,
		)
	}
	return strings.Join(listprime, ","), nil
}

func ComplementCPUList(list string) (string, error) {
	maxcpus, err := TotalAvailable()
	if err != nil {
		return "", err
	}
	return GenerateComplementCPUList(list, maxcpus)
}

// cpuListsExclusive:
// checks if two CPU lists are mutually exclusive
func cpuListsExclusive(list1,
	list2 string, totalCPUs int) (bool, error) {

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

func CPUListsExclusive(list1, list2 string) (bool, error) {
	maxcpus, err := TotalAvailable()
	if err != nil {
		return false, err
	}
	return cpuListsExclusive(list1, list2, maxcpus)
}
