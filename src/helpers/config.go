package helpers

import (
	"fmt"
	"os"

	"github.com/canonical/rt-conf/src/data"
)

func LoadConfigFile(confPath string) (*data.Config, error) {
	_, err := os.Stat(confPath)
	if err != nil {
		return nil, fmt.Errorf("failed to find file: %v", err)
	}

	content, err := ReadYAML(confPath)
	if err != nil {
		return nil, err
	}

	/*
		TODO: Needs to implement proper validation of the parameters
		and parameters format

		validations to be configured:
			- key=value
			- flag
	*/
	return &content, nil
}

// translateConfig translates YAML configuration into kernel command-line parameters.
func TranslateConfig(cfg *data.Config) []string {
	var result []string

	if cfg.KernelCmdline.IsolCPUs != "" {
		result = append(result, Parameters[0].TransformFn(cfg.KernelCmdline.IsolCPUs))
	}

	// TODO: Make this optional
	result = append(result, Parameters[1].TransformFn(cfg.KernelCmdline.DyntickIdle))

	if cfg.KernelCmdline.AdaptiveTicks != "" {
		result = append(result, Parameters[2].TransformFn(cfg.KernelCmdline.AdaptiveTicks))
	}

	return result
}
