package helpers

import (
	"fmt"
	"os"

	"github.com/canonical/rt-conf/src/data"
	"gopkg.in/yaml.v3"
)

func readYAML(path string) (cfg data.Config, err error) {
	d, err := os.ReadFile(path)
	if err != nil {
		// TODO: improve error logging
		fmt.Printf("Failed to read file: %v\n", err)
		return data.Config{}, err
	}

	err = yaml.Unmarshal([]byte(d), &cfg)
	if err != nil {
		// TODO: improve error logging
		fmt.Printf("Failed to unmarshal data: %v\n", err)
		return data.Config{}, err
	}
	return cfg, nil
}
