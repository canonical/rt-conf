package kcmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/canonical/rt-conf/src/model"
)

// UpdateGrub reads GRUB_CMDLINE_LINUX_DEFAULT from the default GRUB configuration file,
// merges it with the kernel command line parameters specified in the provided config,
// and writes the resulting command line to a drop-in configuration file for GRUB.
func UpdateGrub(cfg *model.InternalConfig) ([]string, error) {

	params := model.ConstructKeyValuePairs(&cfg.Data.KernelCmdline)
	if len(params) == 0 {
		return nil, fmt.Errorf("no parameters to inject")
	}

	cmdline, err := parseGrubCMDLineLinuxDefault(cfg.GrubCfg.GrubDefaultFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s grub file: %s",
			cfg.GrubCfg.GrubDefaultFilePath, err.Error())
	}

	if err := duplicatedParams(cmdline); err != nil {
		return nil, fmt.Errorf(
			"invalid existing parameters in %s for GRUB_CMDLINE_LINUX_DEFAULT: %s",
			cfg.GrubCfg.GrubDefaultFilePath, err)
	}
	currParams := model.CmdlineToParams(cmdline)

	// This replaces if the param already exists and
	// creates a new one if it doesn't
	for k, v := range params {
		currParams[k] = v
	}

	cfg.GrubCfg.Cmdline = model.ParamsToCmdline(currParams)
	log.Println("Final kcmdline:", cfg.GrubCfg.Cmdline)

	if err := processFile(cfg.GrubCfg); err != nil {
		return nil, fmt.Errorf("error updating %s: %v",
			cfg.GrubCfg.CustomGrubFilePath, err)
	}

	return GrubConclusion(cfg.GrubCfg.CustomGrubFilePath), nil
}

func parseGrubCMDLineLinuxDefault(path string) (string, error) {
	grubMap, err := ParseDefaultGrubFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to parse grub file: %v", err)
	}
	cmdline, ok := grubMap["GRUB_CMDLINE_LINUX_DEFAULT"]
	if !ok {
		log.Printf("GRUB_CMDLINE_LINUX_DEFAULT not found in %s", path)
	}
	return cmdline, nil
}

func ParseDefaultGrubFile(f string) (map[string]string, error) {
	var err error
	params := make(map[string]string)
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split the line into key and value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		// Trim spaces and quotes from the key and value
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, `"`)

		// Store the key-value pair in the map
		params[key] = value
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading grub file: %v", err)
	}

	return params, err
}

func duplicatedParams(cmdline string) error {
	params := make(map[string]string)
	s := strings.Split(cmdline, " ")
	if len(s) <= 1 {
		// If it's only one parameter, there are no duplicates
		return nil
	}
	for _, p := range s {
		pair := strings.Split(p, "=")
		// Skip parameters without a value, they can be safely dropped
		if len(pair) != 2 {
			// Value is optional for some kernel cmdline parameters
			params[p] = ""
			continue
		}
		param, ok := params[pair[0]]
		if ok {
			// Skip if the value is the same, it can be safelly dropped
			if param == pair[1] {
				continue
			}

			return fmt.Errorf("duplicated parameter: %s=%s and %s=%s",
				pair[0], param, pair[0], pair[1])
		}
		params[pair[0]] = pair[1]
	}
	return nil
}

// processFile processes a file with a given FileTransformer, applying
// its transformation on lines matching the pattern.
var processFile = func(grub model.Grub) error {
	cmdline := "GRUB_CMDLINE_LINUX_DEFAULT=\"" + grub.Cmdline + "\""

	if err := os.WriteFile(grub.CustomGrubFilePath, []byte(cmdline), 0644); err != nil {
		return fmt.Errorf("failed to write to %s file: %v",
			grub.CustomGrubFilePath, err)
	}
	return nil
}
