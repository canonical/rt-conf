package model

import (
	"fmt"
	"os"
	"syscall"
	"testing"

	"github.com/canonical/go-snapctl"
	"go.yaml.in/yaml/v4"
)

var expectedPermission os.FileMode = 0o644

var isOwnedByRoot = func(fi os.FileInfo) bool {
	if testing.Testing() {
		return true // Always return true in test mode
	}
	uid := fi.Sys().(*syscall.Stat_t).Uid
	return uid == 0 // Check if the file is owned by root
}

func (c *Config) LoadFromFile(confPath string) error {
	fileInfo, err := os.Stat(confPath)
	if err != nil {
		return fmt.Errorf("failed to find file: %v", err)
	}

	if fileInfo.Mode() != expectedPermission {
		return fmt.Errorf(
			"file %s has invalid permissions: %v, expected permissions %v",
			confPath, fileInfo.Mode(), expectedPermission)
	}

	if !isOwnedByRoot(fileInfo) {
		return fmt.Errorf("file %s is not owned by root", confPath)
	}

	cfg, err := ReadYAML(confPath)
	if err != nil {
		return err
	}
	*c = *cfg

	/*
		TODO: Needs to implement proper validation of the parameters
		and parameters format

		validations to be configured:
			- key=value
			- flag
	*/
	return nil
}

// LoadSnapOptions reads IRQ and CPU governance objects from snap options
// When a value is set, the whole object gets overridden.
func (c *Config) LoadSnapOptions() error {
	value, err := snapctl.Get(
		"kernel-cmdline",
		"irq-tuning",
		"cpu-governance",
	).Document().Run()
	if err != nil {
		return fmt.Errorf("failed to get snap option: %v", err)
	}

	var confOptions Config

	// Unmarshal json using YAML unmarshaler
	// This works because YAML is a superset of JSON
	err = yaml.Unmarshal([]byte(value), &confOptions)
	if err != nil {
		return fmt.Errorf("failed to unmarshal snap options: %v", err)
	}

	// reject kernel command line arguments
	if len(confOptions.KernelCmdline.Parameters) > 0 {
		return fmt.Errorf("kernel-cmdline snap option is not supported, use the config file instead")
	}

	// override full objects
	if len(confOptions.Interrupts) > 0 {
		c.Interrupts = confOptions.Interrupts
	}
	if len(confOptions.CpuGovernance) > 0 {
		c.CpuGovernance = confOptions.CpuGovernance
	}

	err = c.Validate()
	if err != nil {
		return fmt.Errorf("invalid configurations via snap options: %v", err)
	}

	return nil
}
