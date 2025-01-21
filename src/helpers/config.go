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
