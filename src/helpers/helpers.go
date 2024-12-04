package helpers

import (
	"bufio"
	"fmt"
	"os"

	"github.com/canonical/rt-conf/src/data"
	"github.com/canonical/rt-conf/src/models"
	"github.com/canonical/rt-conf/src/validator"
	"gopkg.in/yaml.v3"
)

// ReadConfigFile reads the configuration file and unmarshals its content
// into the InternalConfig struct.
func ReadConfigFile(cfg *InternalConfig) error {
	/*
		TODO: Needs to implement proper validation of the parameters
		and parameters format

		validations to be configured:
			- key=value
			- flag
	*/

	data, err := os.ReadFile(cfg.ConfigFile)
	if err != nil {
		// TODO: improve error logging
		fmt.Printf("Failed to read file: %v\n", err)
		return err
	}

	err = yaml.Unmarshal([]byte(data), &cfg.Data)
	if err != nil {
		// TODO: improve error logging
		fmt.Printf("Failed to unmarshal data: %v\n", err)
		return err
	}
	return nil
}

// processFile processes a file with a given FileTransformer, applying
// its transformation on lines matching the pattern.
func ProcessFile(transformer data.FileTransformer) error {
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

type InternalConfig struct {
	ConfigFile string
	Data       data.Config

	GrubDefault data.Grub

	// GrubCfg data.Grub
}

var Parameters = []data.Param{
	{
		YAMLName:    "isolcpus",
		CmdlineName: "isolcpus",
		TransformFn: func(value interface{}) string {
			return fmt.Sprintf("isolcpus=%s", value)
		},
	},
	{
		YAMLName:    "dyntick-idle",
		CmdlineName: "nohz",
		TransformFn: func(value interface{}) string {
			validator.ValidateType(validator.TypeEnum["bool"], "dyntick-idle",
				value)
			if v, ok := value.(bool); ok && v {
				return "nohz=on"
			}
			return "nohz=off"
		},
	},
	{
		YAMLName:    "adaptive-ticks",
		CmdlineName: "nohz_full",
		TransformFn: func(value interface{}) string {
			return fmt.Sprintf("nohz_full=%s", value)
		},
	},
}

// translateConfig translates YAML configuration into kernel command-line parameters.
func translateConfig(cfg data.Config) []string {
	var result []string
	for _, param := range Parameters {
		if value, exists := cfg[param.YAMLName]; exists {
			result = append(result, param.TransformFn(value))
		}
	}
	return result
}

// InjectToGrubFiles inject the kernel command line parameters to the grub files.
// /boot/grub/grub.cfg and /etc/default/grub
func (c *InternalConfig) InjectToGrubFiles() error {
	err := ReadConfigFile(c)
	if err != nil {
		return err
	}

	kernelCmdline := c.Data["kernel-cmdline"].(data.Config)
	cmdline := translateConfig(kernelCmdline)
	fmt.Println("KernelCmdline: ", cmdline)

	// rpiCmdline := &rpiTransformer{
	// 	filePath: c.grubCfg.file,
	// 	pattern:  c.grubCfg.pattern,
	// 	params:   c.data.KernelCmdline,
	// }

	// grubCfg := &models.GrubCfgTransformer{
	// 	FilePath: c.GrubCfg.File,
	// 	Pattern:  c.GrubCfg.Pattern,
	// 	Params:   cmdline,
	// }

	grubDefault := &models.GrubDefaultTransformer{
		FilePath: c.GrubDefault.File,
		Pattern:  c.GrubDefault.Pattern,
		Params:   cmdline,
	}

	if err := ProcessFile(grubDefault); err != nil {
		return fmt.Errorf("error updating %s: %v", grubDefault.FilePath, err)
	}
	fmt.Printf("File %v updated successfully.\n", grubDefault.FilePath)

	// // Process each file with its specific transformer
	// if err := ProcessFile(grubCfg); err != nil {
	// 	return fmt.Errorf("error updating %s: %v", grubCfg.FilePath, err)
	// }
	// fmt.Println("File /boot/grub/grub.cfg updated successfully.")

	return nil
}
