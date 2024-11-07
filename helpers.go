package main

import (
	"bufio"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ReadConfigFile reads the configuration file and unmarshals its content
// into the InternalConfig struct.
func readConfigFile(cfg *InternalConfig) error {
	/*
		TODO: Needs to implement proper validation of the parameters
		and parameters format

		validations to be configured:
			- key=value
			- flag
	*/

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

// processFile processes a file with a given FileTransformer, applying
// its transformation on lines matching the pattern.
func processFile(transformer FileTransformer) error {
	file, err := os.Open(transformer.GetFilePath())
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Create a temporary file to write modified content
	tmpFile, err := os.CreateTemp("", "config-modified-")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if transformer.GetPattern().MatchString(line) {
			line = transformer.TransformLine(line)
		}
		_, err := tmpFile.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("failed to write to temp file: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// Replace the original file with the modified one
	tmpFile.Close()
	err = os.Rename(tmpFile.Name(), transformer.GetFilePath())
	if err != nil {
		return fmt.Errorf("failed to replace original file: %v", err)
	}
	return nil
}
