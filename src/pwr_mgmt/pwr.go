package pwrmgmt

import (
	"fmt"
	"log"
	"os"

	"github.com/canonical/rt-conf/src/cpulists"
	"github.com/canonical/rt-conf/src/model"
)

type ReaderWriter struct {
	Path string
}

var scalingGovernorReaderWriter = ReaderWriter{
	Path: "/sys/devices/system/cpu/cpu%d/cpufreq/scaling_governor",
}

func (w ReaderWriter) WriteScalingGov(sclgov string, cpu int) error {
	scalingGovFile := fmt.Sprintf(w.Path, cpu)

	err := os.WriteFile(scalingGovFile, []byte(sclgov), 0644)
	if err != nil {
		return fmt.Errorf("error writing to %s: %v", scalingGovFile, err)
	} else {
		log.Printf("Set %s to %s", scalingGovFile, sclgov)
	}

	return nil
}

func ApplyPwrConfig(config *model.InternalConfig) error {
	log.Println("\n-----------------------")
	log.Println("Applying CPU Governance")
	log.Println("-----------------------")

	if len(config.Data.CpuGovernance) == 0 {
		log.Println("No CPU governance rules found in config")
		return nil
	}
	return scalingGovernorReaderWriter.applyPwrConfig(config.Data.CpuGovernance)
}

// Apply changes based on YAML config
func (wr ReaderWriter) applyPwrConfig(
	config []model.CpuGovernanceRule,
) error {

	// Range over all CPU governance rules
	for _, sclgov := range config {
		cpus, err := cpulists.Parse(sclgov.CPUs)
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
