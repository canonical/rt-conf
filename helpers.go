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
	// Open file with read and write permissions
	file, err := os.OpenFile(transformer.GetFilePath(), os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Read all lines into a slice
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if transformer.GetPattern().MatchString(line) {
			line = transformer.TransformLine(line)
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// Truncate file and write transformed lines
	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate file: %v", err)
	}
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek to start of file: %v", err)
	}

	for _, line := range lines {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("failed to write to file: %v", err)
		}
	}

	return nil
}
