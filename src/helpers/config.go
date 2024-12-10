package helpers

import (
	"fmt"
	"os"

	"github.com/canonical/rt-conf/src/data"
)

func LoadConfigFile(confPath, grubPath string) (data.InternalConfig, error) {
	_, err := os.Stat(confPath)
	if err != nil {
		return data.InternalConfig{}, fmt.Errorf("failed to find file: %v", err)
	}

	content, err := ReadYAML(confPath)
	if err != nil {
		return data.InternalConfig{}, err
	}

	/*
		TODO: Needs to implement proper validation of the parameters
		and parameters format

		validations to be configured:
			- key=value
			- flag
	*/
	return data.InternalConfig{
		Data: content,
		GrubDefault: data.Grub{
			File:    grubPath,
			Pattern: data.PatternGrubDefault,
		},
	}, nil
}

// translateConfig translates YAML configuration into kernel command-line parameters.
func TranslateConfig(cfg data.Config) []string {
	var result []string
	result = append(result, Parameters[0].TransformFn(cfg.KernelCmdline.IsolCPUs))
	result = append(result, Parameters[1].TransformFn(cfg.KernelCmdline.DyntickIdle))
	result = append(result, Parameters[2].TransformFn(cfg.KernelCmdline.AdaptiveTicks))
	return result
}
