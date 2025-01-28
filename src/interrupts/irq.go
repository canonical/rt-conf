package interrupts

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/canonical/rt-conf/src/common"
	"github.com/canonical/rt-conf/src/cpu"
	"github.com/canonical/rt-conf/src/data"
)

// NOTE: to be able to remap IRQs:
// CONFIG_REGMAP_IRQ=y must be present in the kernel config
// cat /boot/config-`uname -r` | grep CONFIG_REGMAP_IRQ

// /proc/interrupts man page:
// https://man7.org/linux/man-pages/man5/proc_interrupts.5.html

// OBS: There is no man page for /proc/irq

// NOTE: The Documentation for the procfs is available at:
// https://docs.kernel.org/filesystems/proc.html

// The Kernel space API irq_create_mapping_affinity() returns the IRQ number
// for a given interrupt line and affinity mask. This is used to map the
// interrupt line to the IRQ number.
// https://elixir.bootlin.com/linux/v6.13-rc2/source/kernel/irq/irqdomain.c#L834

// About SMI (System Management Interrupt):
// https://wiki.linuxfoundation.org/realtime/documentation/howto/debugging/smi-latency/smi

// RealIRQWriter writes CPU affinity to the real `/proc/irq/<irq>/smp_affinity_list` file.
type RealIRQWriter struct{}

// RealIRQReader reads IRQs from the real `/sys/kernel/irq` directory.
type RealIRQReader struct{}

// Write IRQ affinity
func (w *RealIRQWriter) WriteCPUAffinity(irqNum, cpus string) error {
	affinityFile :=
		fmt.Sprintf("%s/%s/smp_affinity_list", common.ProcIRQ, irqNum)
	return os.WriteFile(affinityFile, []byte(cpus), 0644)
}

func (r *RealIRQReader) ReadIRQs() ([]IRQInfo, error) {
	var irqInfos []IRQInfo

	// Read the directories in /sys/kernel/irq
	dirEntries, err := os.ReadDir(common.SysKernelIRQ)
	if err != nil {
		return nil, err
	}

	for _, entry := range dirEntries {
		if entry.IsDir() {
			number, err := strconv.ParseUint(entry.Name(), 10, 32)
			if err != nil {
				/* This may happen if the kernel IRQ structure
				evolves sometime or somehow in the future */
				continue // Skip if not a valid number
			}
			var irqInfo IRQInfo
			irqInfo.Number = uint(number)

			// Read files in the IRQ directory
			files := []string{
				"actions", "chip_name", "name", "type", "wakeup",
			}
			for _, file := range files {
				filePath := filepath.Join(
					common.SysKernelIRQ, entry.Name(), file,
				)
				content, err := os.ReadFile(filePath)
				if err != nil {
					// TODO: Log warning here
					continue
				}
				c := strings.TrimSuffix(
					strings.TrimSpace(string(content)), "\n")
				switch file {
				case "actions":
					irqInfo.Actions = c
				case "chip_name":
					irqInfo.ChipName = c
				case "name":
					irqInfo.Name = c
				case "type":
					irqInfo.Type = c
				case "wakeup":
					irqInfo.Wakeup = c
				}
			}
			irqInfos = append(irqInfos, irqInfo)
		}
	}
	return irqInfos, err
}

// Apply changes based on YAML config
func ApplyIRQConfig(
	config *data.InternalConfig,
	reader IRQReader,
	writer IRQWriter,
) error {

	irqs, err := reader.ReadIRQs()
	if err != nil {
		return err
	}

	// Range over IRQ tunning array
	for _, irqTuning := range config.Data.Interrupts {
		matchingIRQs, err := filterIRQs(irqs, irqTuning.Filter)
		if err != nil {
			return err
		}

		for _, irqNum := range matchingIRQs {
			err := writer.WriteCPUAffinity(irqNum, irqTuning.CPUs)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// TODO: Refact this function to shrink the size
// Filter IRQs by criteria
func filterIRQs(
	irqs []IRQInfo,
	filter data.IRQFilter) ([]string, error) {
	var matchingIRQs []string

	for _, entry := range irqs {
		if !strings.HasPrefix(entry.Name, "irq") {
			continue
		}
		irqNum := entry.Name
		irqPath := filepath.Join(common.SysKernelIRQ, irqNum)

		// Apply filters
		match, err := matchFilter(filepath.Base(irqPath), filter.Number)
		if err != nil {
			return nil, err
		}
		if filter.Number != "" && !match {
			continue
		}

		match, err = matchFile(filepath.Join(irqPath, "smp_affinity_list"),
			filter.Action)
		if err != nil {
			return nil, err
		}
		if filter.Action != "" && !match {
			continue
		}

		match, err = matchFile(filepath.Join(irqPath, "chip_name"),
			filter.ChipName)
		if err != nil {
			return nil, err
		}
		if filter.Action != "" && !match {
			continue
		}

		match, err = matchFile(filepath.Join(irqPath, "name"), filter.Name)
		if err != nil {
			return nil, err
		}
		if filter.ChipName != "" && !match {
			continue
		}

		match, err = matchFile(filepath.Join(irqPath, "type"), filter.Type)
		if err != nil {
			return nil, err
		}
		if filter.Name != "" && !match {
			continue
		}

		match, err = matchFile(filepath.Join(irqPath, "type"), filter.Type)
		if err != nil {
			return nil, err
		}
		if filter.Type != "" && !match {
			continue
		}

		matchingIRQs = append(matchingIRQs, irqNum)
	}
	return matchingIRQs, nil
}

func matchFilter(irqPath, filter string) (bool, error) {
	// NOTE: The here the filter is already valited as a cpulist
	irqs, err := cpu.ParseCPUs(filter)
	if err != nil {
		return false, err
	}
	for irq := range irqs {
		if strings.Contains(irqPath, strconv.Itoa(irq)) {
			return true, nil
		}
	}
	return false, nil
}

// Match criteria in a file
func matchFile(filePath, pattern string) (bool, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false, err
	}
	match, err := regexp.MatchString(pattern, string(content))
	if err != nil {
		return false, err
	}
	return match, nil
}
