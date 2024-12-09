package helpers

import (
	"fmt"
	"os"
	"regexp"

	"github.com/canonical/rt-conf/src/data"
)

func LoadConfigFile(confPath, grubPath string) (InternalConfig, error) {
	var content interface{}
	for _, p := range []string{confPath, grubPath} {
		if _, err := os.Stat(p); err != nil {
			return InternalConfig{}, fmt.Errorf("failed to find file: %v", err)
		}
	}

	content, err := readYAML(confPath)
	if err != nil {
		return InternalConfig{}, err
	}

	/*
		TODO: Needs to implement proper validation of the parameters
		and parameters format

		validations to be configured:
			- key=value
			- flag
	*/
	return InternalConfig{
		ConfigFile: confPath,
		Data:       content.(data.Config),
		GrubDefault: data.Grub{
			File:    grubPath,
			Pattern: regexp.MustCompile(data.RegexGrubDefault),
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
