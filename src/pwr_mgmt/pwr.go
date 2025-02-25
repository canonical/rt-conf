package pwrmgmt

import (
	"fmt"
	"log"
	"os"

	"github.com/canonical/rt-conf/src/cpu"
	"github.com/canonical/rt-conf/src/data"
)

type ReaderWriter struct {
	Path string
}

var scalingGovernerReaderWriter = ReaderWriter{
	Path: "/sys/devices/system/cpu/cpu%d/cpufreq/scaling_governor",
}

func (w ReaderWriter) WriteScalingGov(sclgov string, cpu int) error {
	scallingGovFile := fmt.Sprintf(w.Path, cpu)

	err := os.WriteFile(scallingGovFile, []byte(sclgov), 0644)
	if err != nil {
		return fmt.Errorf("error writing to %s: %v", scallingGovFile, err)
	} else {
		log.Printf("Set %s to %s", scallingGovFile, sclgov)
	}

	return nil
}

func ApplyPwrConfig(config *data.InternalConfig) error {
	return scalingGovernerReaderWriter.applyPwrConfig(config.Data.CpuGovernance)
}

// Apply changes based on YAML config
func (wr ReaderWriter) applyPwrConfig(
	config []data.CpuGovernanceRule,
) error {

	// Range over all CPU governance rules
	for _, sclgov := range config {
		cpus, err := cpu.ParseCPUs(sclgov.CPUs)
		if err != nil {
			return err
		}
		for cpu := range cpus {
			err := wr.WriteScalingGov(sclgov.ScalGov, cpu)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
