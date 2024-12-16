package interrupts

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/canonical/rt-conf/src/cpu"
	"github.com/canonical/rt-conf/src/data"
	"github.com/canonical/rt-conf/src/helpers"
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

type ProcInterrupts struct {
	Number int
	CPUs   cpu.CPUs
	Name   string
}

func ProcessIRQIsolation(cfg *data.InternalConfig) error {
	maxcpus, err := cpu.TotalAvailable()
	if err != nil {
		return fmt.Errorf("error getting total CPUs: %v", err)
	}

	irqs, err := systemIRQs()
	if err != nil {
		return fmt.Errorf("error mapping IRQs: %v", err)
	}

	isolCPUs := cfg.Data.Interrupts.IsolateCPU
	newAffinity := cfg.Data.Interrupts.IRQHandler
	if newAffinity == "" {
		var err error
		newAffinity, err = cpu.GenerateComplementCPUList(isolCPUs, maxcpus)
		if err != nil {
			return fmt.Errorf("error generating complement CPU list: %v", err)
		}
	} else {
		excl, err := cpu.MutuallyExclusive(isolCPUs, newAffinity, maxcpus)
		if err != nil {
			return fmt.Errorf("error checking cpu list mutual exclusion: %v", err)
		}
		if !excl {
			return fmt.Errorf("invalid input: cpu lists not mutually excluded: '%v', '%v'", isolCPUs, newAffinity)
		}
	}

	if err := remapIRQsAffinity(newAffinity, irqs); err != nil {
		return fmt.Errorf("error performing CPU isolation: %v", err)
	}

	return nil
}

func remapIRQsAffinity(newAffinity string, irq []uint) error {
	maxcpus, err := cpu.TotalAvailable()
	if err != nil {
		return fmt.Errorf("error getting total CPUs: %v", err)
	}
	fmt.Println("Total CPUs:", maxcpus)
	for _, i := range irq {
		f := fmt.Sprintf("/proc/irq/%d/smp_affinity_list", i)
		err := helpers.WriteToFile(f, newAffinity)
		if err == nil {
			continue
		}
		// SMI IRQs are not allowed to be written to from userspace.
		// It fails with "input/output error"
		if !strings.Contains(err.Error(), "input/output error") {
			return fmt.Errorf("error writing to %s: %v", f, err)
		}
		log.Printf("Managed IRQ %d, skipped\n", i)
	}
	return nil
}

// Get IRQs from /proc/irq
func systemIRQs() ([]uint, error) {
	dirEntries, err := os.ReadDir("/proc/irq")
	if err != nil {
		return nil, fmt.Errorf("error reading /proc/irq directory: %v", err)
	}

	var irqNumbers []uint
	for _, entry := range dirEntries {
		if !entry.IsDir() {
			continue
		}

		irqNumber, err := strconv.Atoi(entry.Name())
		if err != nil {
			return nil, fmt.Errorf("error converting %s to int: %v",
				entry.Name(), err)
		}

		irqNumbers = append(irqNumbers, uint(irqNumber))
	}

	return irqNumbers, nil
}
