package interrupts

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/canonical/rt-conf/src/data"
)

// TODO: THis needs to be superseed in the unit tests
var prefixPath = ""

var filterPath = prefixPath + "/sys/kernel/irq"
var applyPath = prefixPath + "/proc/irq"

// Filter IRQs by criteria
func filterIRQs(basePath string, filter data.IRQFilter) ([]string, error) {
	var matchingIRQs []string

	entries, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "irq") {
			continue
		}
		irqNum := entry.Name()
		irqPath := filepath.Join(basePath, irqNum)

		// Apply filters

		if filter.Number != "" && !matchFilter(filepath.Base(irqPath), filter.Number) {
			continue
		}

		if filter.Action != "" && !matchFile(filepath.Join(irqPath, "actions"), filter.Action) {
			continue
		}
		if filter.ChipName != "" && !matchFile(filepath.Join(irqPath, "chip_name"), filter.ChipName) {
			continue
		}
		if filter.Name != "" && !matchFile(filepath.Join(irqPath, "name"), filter.Name) {
			continue
		}
		if filter.Type != "" && !matchFile(filepath.Join(irqPath, "type"), filter.Type) {
			continue
		}

		matchingIRQs = append(matchingIRQs, irqNum)
	}
	return matchingIRQs, nil
}

// TODO: implement this
func matchFilter(irqPath, filter string) bool {
	return true
}

// Match criteria in a file
func matchFile(filePath, pattern string) bool {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}
	return strings.Contains(string(content), pattern)
}

// Write IRQ affinity
func applyCPUAffinity(irqNum, cpus string) error {
	affinityFile := fmt.Sprintf("%s/proc/irq/%s/smp_affinity_list", prefixPath,
		irqNum)
	return os.WriteFile(affinityFile, []byte(cpus), 0644)
}

// Apply changes based on YAML config
func applyConfig(config *data.Config) error {
	for _, irqTuning := range config.Interrupts {
		matchingIRQs, err := filterIRQs(applyPath, irqTuning.Filter)
		if err != nil {
			return err
		}
		for _, irqNum := range matchingIRQs {
			err := applyCPUAffinity(irqNum, irqTuning.CPUs)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
