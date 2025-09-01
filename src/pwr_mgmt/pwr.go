package pwrmgmt

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/canonical/rt-conf/src/cpulists"
	"github.com/canonical/rt-conf/src/model"
	"github.com/canonical/rt-conf/src/utils"
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

func writeOnly(path string, data string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		return fmt.Errorf("error opening %s: %v", path, err)
	}
	defer func(f *os.File) error {
		if err := f.Close(); err != nil {
			return err
		}
		return nil
	}(f)

	_, err = f.Write([]byte(data))
	if err != nil {
		return fmt.Errorf("error writing to %s: %v", path, err)
	}
	return nil
}

func (w ReaderWriter) WriteScalingGov(sclgov string, cpu int) error {
	if sclgov == "" {
		return nil // No scaling governor set, nothing to write
	}
	scalingGovFile := fmt.Sprintf(w.ScalingGovernorPath, cpu)

	err := writeOnly(scalingGovFile, sclgov)
	if err != nil {
		return fmt.Errorf("error writing to %s: %v", scalingGovFile, err)
	}
	return nil
}

func (w ReaderWriter) WriteCPUFreq(freqMin, freqMax, cpu int) error {
	if freqMin != -1 {
		minFreqSysfs := fmt.Sprintf(w.MinFreqPath, cpu)
		if err := writeOnly(minFreqSysfs,
			strconv.Itoa(freqMin)); err != nil {
			return fmt.Errorf("error writing to %s: %v", minFreqSysfs, err)
		}
	}

	if freqMax != -1 {
		maxFreqSysfs := fmt.Sprintf(w.MaxFreqPath, cpu)
		if err := writeOnly(maxFreqSysfs,
			strconv.Itoa(freqMax)); err != nil {
			return fmt.Errorf("error writing to %s: %v", maxFreqSysfs, err)
		}
	}

	return nil
}

func ApplyPwrConfig(config *model.InternalConfig) error {
	utils.PrintTitle("CPU Governance")
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

		log.Printf("Rule: %s \n", label)
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
	pluralSuffix := "s"
	if len(cpus) == 1 {
		pluralSuffix = ""
	}
	cpuList := cpulists.GenCPUlist(cpus)

	var msg []string
	if scalingGov != "" {
		msg = append(msg,
			fmt.Sprintf("Set scaling governance of CPU%s %s to %s", pluralSuffix,
				cpuList, scalingGov))
	}
	if minFreq != "" {
		msg = append(msg,
			fmt.Sprintf("Set min frequency of CPU%s %s to %s", pluralSuffix,
				cpuList, minFreq))
	}
	if maxFreq != "" {
		msg = append(msg,
			fmt.Sprintf("Set max frequency of CPU%s %s to %s", pluralSuffix,
				cpuList, maxFreq))
	}

	utils.LogTreeStyle(msg)
}

func (wr ReaderWriter) applyRule(cpu int, sclgov model.CpuGovernanceRule) error {
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
