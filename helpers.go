package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func readConfigFile(cfg *InternalConfig) error {

	data, err := os.ReadFile(cfg.configFile)
	if err != nil {
		// TODO: improve error logging
		fmt.Printf("Failed to read file: %v\n", err)
		return err
	}

	err = yaml.Unmarshal([]byte(data), &cfg.data)
	if err != nil {
		// TODO: improve error logging
		fmt.Printf("Failed to unmarshal data: %v\n", err)
		return err
	}
	return nil
}
