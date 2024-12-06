package helpers

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/canonical/rt-conf/src/data"
	"github.com/canonical/rt-conf/src/execute"
	"github.com/canonical/rt-conf/src/models"
	"github.com/canonical/rt-conf/src/validator"
	"gopkg.in/yaml.v3"
)

type InternalConfig struct {
	ConfigFile string
	Data       data.Config

	GrubDefault data.Grub
}

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
			Pattern: regexp.MustCompile(models.RegexGrubDefault),
		},
	}, nil

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
func TranslateConfig(cfg data.Config) []string {
	var result []string
	result = append(result, Parameters[0].TransformFn(cfg.KernelCmdline.IsolCPUs))
	result = append(result, Parameters[1].TransformFn(cfg.KernelCmdline.DyntickIdle))
	result = append(result, Parameters[2].TransformFn(cfg.KernelCmdline.AdaptiveTicks))
	return result
}

// InjectToGrubFiles inject the kernel command line parameters to the grub files.
// /boot/grub/grub.cfg and /etc/default/grub
func UpdateGrub(cfg *InternalConfig) error {
	cmdline := translateConfig(cfg.Data)
	fmt.Println("KernelCmdline: ", cmdline)

	grubDefault := &models.GrubDefaultTransformer{
		FilePath: cfg.GrubDefault.File,
		Pattern:  cfg.GrubDefault.Pattern,
		Params:   cmdline,
	}

	if err := ProcessFile(grubDefault); err != nil {
		return fmt.Errorf("error updating %s: %v", grubDefault.FilePath, err)
	}
	fmt.Printf("File %v updated successfully.\n", grubDefault.FilePath)

	execute.GrubConclusion()
	return nil
}
