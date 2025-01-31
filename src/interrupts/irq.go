package interrupts

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/canonical/rt-conf/src/common"
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

// ** From experiments:
// ** Non active IRQs (not shown in /proc/interrupts) are the ones which
// ** doesn't have an action (/sys/kernel/irq/<num>/action) associated with them

type IRQs map[int]bool // use the same logic as CPUs lists

// realIRQReaderWriter writes CPU affinity to the real `/proc/irq/<irq>/smp_affinity_list` file.
type realIRQReaderWriter struct{}

// Write IRQ affinity
func (w *realIRQReaderWriter) WriteCPUAffinity(irqNum int, cpus string) error {
	affinityFile :=
		fmt.Sprintf("%s/%d/smp_affinity_list", common.ProcIRQ, irqNum)

	err := os.WriteFile(affinityFile, []byte(cpus), 0644)
	// SMI are not allowed to be written to from userspace.
	// It fails with "input/output error" this error can be ignored.
	if err != nil {
		if !strings.Contains(err.Error(), "input/output error") {
			return fmt.Errorf("error writing to %s: %v", affinityFile, err)
		}
	}
	return nil
}

func (r *realIRQReaderWriter) ReadIRQs() ([]IRQInfo, error) {
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
			irqInfo.Number = int(number)

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

func ApplyIRQConfig(config *data.InternalConfig) error {
	return applyIRQConfig(config, &realIRQReaderWriter{})
}

// Apply changes based on YAML config
func applyIRQConfig(
	config *data.InternalConfig,
	handler IRQReaderWriter,
) error {

	irqs, err := handler.ReadIRQs()
	if err != nil {
		return err
	}

	if len(irqs) == 0 {
		return fmt.Errorf("no IRQs found")
	}

	// Range over IRQ tunning array
	for _, irqTuning := range config.Data.Interrupts {
		matchingIRQs, err := filterIRQs(irqs, irqTuning.Filter)
		if err != nil {
			return fmt.Errorf("failed to filter IRQs: %v", err)
		}

		if len(matchingIRQs) == 0 {
			//TODO: log warning here
			// TODO: confirm if it should fail when nothing is matched
			return fmt.Errorf("no IRQs matched the filter: %v",
				irqTuning.Filter)
		}

		for irqNum := range matchingIRQs {
			err := handler.WriteCPUAffinity(irqNum, irqTuning.CPUs)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// filterIRQs filters IRQs based on the provided filters (matches any filter).
func filterIRQs(irqs []IRQInfo, filter data.IRQFilter) (IRQs, error) {
	matchingIRQs := make(IRQs)

	for _, irq := range irqs {
		if matchesAnyFilter(irq, filter) {
			matchingIRQs[irq.Number] = true
		}
	}
	return matchingIRQs, nil
}

// matchesAnyFilter checks if an IRQ matches any of the given filters.
func matchesAnyFilter(irq IRQInfo, filter data.IRQFilter) bool {
	// TODO: these needs to be AND
	return matchesRegex(irq.Actions, filter.Action) &&
		matchesRegex(irq.ChipName, filter.ChipName) &&
		matchesRegex(irq.Name, filter.Name) &&
		matchesRegex(irq.Type, filter.Type)
}

// matchesRegex checks if a field matches a regex pattern.
func matchesRegex(value, pattern string) bool {
	if pattern == "" {
		// TODO: if empty return true since it's an AND
		return true
	}
	match, err := regexp.MatchString(pattern, value)
	return err == nil && match
}
