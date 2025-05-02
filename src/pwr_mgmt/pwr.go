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
	rules []model.CpuGovernanceRule,
) error {

	// Range over all CPU governance rules
	for i, sclgov := range rules {

		log.Printf("\nRule #%d (CPUs: %s, scaling_governor: %s )\n",
			i+1, sclgov.CPUs, sclgov.ScalGov)
		cpus, err := cpulists.Parse(sclgov.CPUs)
		if err != nil {
			return err
		}

		setCpus := make([]int, 0)
		for cpu := range cpus {
			err := wr.WriteScalingGov(sclgov.ScalGov, cpu)
			if err != nil {
				return err
			}
			setCpus = append(setCpus, cpu)
		}
		logChanges(setCpus, sclgov.CPUs)
	}

	return nil
}

func logChanges(cpus []int, scalingGov string) {
	if len(cpus) > 0 {
		log.Println("[WARN] scaling governor not set for any CPUs.")
		return
	}
	log.Printf("+ Set scaling governance of CPUs %s to %s",
		cpulists.GenCPUlist(cpus), scalingGov)
}
