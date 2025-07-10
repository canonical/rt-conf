package pwrmgmt

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/canonical/rt-conf/src/cpulists"
	"github.com/canonical/rt-conf/src/model"
)

type ReaderWriter struct {
	ScalingGovernorPath string
	MinFreqPath         string
	MaxFreqPath         string
}

var pwrmgmtReaderWriter = ReaderWriter{
	ScalingGovernorPath: "/sys/devices/system/cpu/cpu%d/cpufreq/scaling_governor",
	MinFreqPath:         "/sys/devices/system/cpu/cpu%d/cpufreq/scaling_min_freq",
	MaxFreqPath:         "/sys/devices/system/cpu/cpu%d/cpufreq/scaling_max_freq",
}

func (w ReaderWriter) WriteScalingGov(sclgov string, cpu int) error {
	if sclgov == "" {
		return nil // No scaling governor set, nothing to write
	}
	scalingGovFile := fmt.Sprintf(w.ScalingGovernorPath, cpu)

	err := os.WriteFile(scalingGovFile, []byte(sclgov), 0644)
	if err != nil {
		return fmt.Errorf("error writing to %s: %v", scalingGovFile, err)
	}
	return nil
}

func (w ReaderWriter) WriteCPUFreq(freqMin, freqMax, cpu int) error {

	if freqMin != -1 {
		minFreqSysfs := fmt.Sprintf(w.MinFreqPath, cpu)
		if err := os.WriteFile(minFreqSysfs, []byte(strconv.Itoa(freqMin)),
			0644); err != nil {
			return fmt.Errorf("error writing to %s: %v", minFreqSysfs, err)
		}
	}

	if freqMax != -1 {
		maxFreqSysfs := fmt.Sprintf(w.MaxFreqPath, cpu)
		if err := os.WriteFile(maxFreqSysfs, []byte(strconv.Itoa(freqMax)),
			0644); err != nil {
			return fmt.Errorf("error writing to %s: %v", maxFreqSysfs, err)
		}
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
	return pwrmgmtReaderWriter.applyPwrConfig(config.Data.CpuGovernance)
}

// Apply changes based on YAML config
func (wr ReaderWriter) applyPwrConfig(
	rules model.PwrMgmt,
) error {

	// Range over all CPU governance rules
	for label, sclgov := range rules {

		log.Printf("\nRule: %s \n", label)
		cpus, err := cpulists.Parse(sclgov.CPUs)
		if err != nil {
			return err
		}

		var setCpus []int

		for cpu := range cpus {
			if err := wr.applyRule(cpu, sclgov); err != nil {
				return fmt.Errorf("failed to apply CPU governance rule #%s for CPU %d: %v",
					label, cpu, err)
			}
			setCpus = append(setCpus, cpu)
		}
		logChanges(setCpus, sclgov.MinFreq, sclgov.MaxFreq, sclgov.ScalGov)
	}

	return nil
}

func logChanges(cpus []int, minFreq, maxFreq, scalingGov string) {
	cpusName := "CPUs"
	if len(cpus) == 1 {
		cpusName = "CPU"
	}
	log.Printf("+ Set scaling governance of %s %s to %s\n",
		cpusName, cpulists.GenCPUlist(cpus), scalingGov)

	minFreqConnector := "└── "
	if minFreq != "" {
		if maxFreq != "" {
			minFreqConnector = "├── "
		}
		log.Printf("%sSet min frequency of %s %s to %s\n",
			minFreqConnector, cpusName, cpulists.GenCPUlist(cpus), minFreq)
	}
	if maxFreq != "" {
		log.Printf("└── Set max frequency of %s %s to %s\n",
			cpusName, cpulists.GenCPUlist(cpus), maxFreq)
	}
}

func (wr ReaderWriter) applyRule(cpu int,
	sclgov model.CpuGovernanceRule) error {
	if err := wr.WriteScalingGov(sclgov.ScalGov, cpu); err != nil {
		return err
	}
	minFreq, err := model.ParseFreq(sclgov.MinFreq)
	if err != nil {
		return err
	}
	maxFreq, err := model.ParseFreq(sclgov.MaxFreq)
	if err != nil {
		return err
	}
	if err := wr.WriteCPUFreq(
		minFreq,
		maxFreq,
		cpu); err != nil {
		return fmt.Errorf("failed to set CPU frequency for CPU %d: %v", cpu, err)
	}
	return nil
}
