package helpers

import (
	"fmt"
	"os"

	"github.com/canonical/rt-conf/src/data"
	"gopkg.in/yaml.v3"
)

func ReadYAML(path string) (cfg *data.Config, err error) {
	d, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	err = yaml.Unmarshal([]byte(d), &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %v", err)
	}

	err = cfg.Validate()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
