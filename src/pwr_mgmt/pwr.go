package pwrmgmt

import (
	"fmt"
	"log"
	"os"

	"github.com/canonical/rt-conf/src/cpu"
	"github.com/canonical/rt-conf/src/data"
)

const scalingGovernorPath = "/sys/devices/system/cpu/cpu%d/cpufreq/scaling_governor"

// realScalGovReaderWriter writes CPU scalling governor string to
// the real `/sys/devices/system/cpu/cpu%d/cpufreq/scaling_governor` file.
type realScalGovReaderWriter struct{}

func (w *realScalGovReaderWriter) WriteScalingGov(sclgov string, cpu int) error {
	scallingGovFile :=
		fmt.Sprintf(scalingGovernorPath, cpu)

	err := os.WriteFile(scallingGovFile, []byte(sclgov), 0644)
	if err != nil {
		return fmt.Errorf("error writing to %s: %v", scallingGovFile, err)
	} else {
		log.Printf("Set %s to %s", scallingGovFile, sclgov)
	}

	return nil
}

func (r *realScalGovReaderWriter) ReadPwrSetting() ([]PwrInfo, error) {
	return nil, nil
}

func ApplyPwrConfig(config *data.InternalConfig) error {
	return applyPrwConfig(config, &realScalGovReaderWriter{})
}

// Apply changes based on YAML config
func applyPrwConfig(
	config *data.InternalConfig,
	handler ScalGovReaderWriter,
) error {

	// Range over all CPU governance rules
	for _, sclgov := range config.Data.CpuGovernance {
		cpus, _ := cpu.ParseCPUs(sclgov.CPUs)
		for cpu := range cpus {
			err := handler.WriteScalingGov(sclgov.ScalGov, cpu)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
