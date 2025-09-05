package model

import (
	"fmt"
	"os"
	"regexp"

	"go.yaml.in/yaml/v4"
)

type Grub struct {
	GrubDropInFile string
	Cmdline        string
}

type Core interface {
	InjectToFile(pattern *regexp.Regexp) error
}

func ReadYAML(path string) (cfg *Config, err error) {
	d, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	err = yaml.Unmarshal([]byte(d), &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %v", err)
	}
	if cfg == nil {
		return nil, fmt.Errorf("empty config file")
	}
	err = cfg.Validate()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
