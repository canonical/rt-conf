package data

import (
	"regexp"
)

// FileTransformer interface with a TransformLine method.
// This method is used to transform a line of a file.
//
// NOTE: This interface can be implemented also for RPi on classic
type FileTransformer interface {
	TransformLine(string) string
	GetFilePath() string
	GetPattern() *regexp.Regexp
}

type Core interface {
	InjectToFile(pattern *regexp.Regexp) error
}

type Grub struct {
	File    string
	Pattern *regexp.Regexp
}

type Config struct {
	KernelCmdline KernelCmdline `yaml:"kernel-cmdline"`
}

// KernelCmdline represents the kernel command line options.
type KernelCmdline struct {
	IsolCPUs      string `yaml:"isolcpus"`       // Isolate CPUs
	DyntickIdle   bool   `yaml:"dyntick-idle"`   // Enable/Disable dyntick idle
	AdaptiveTicks string `yaml:"adaptive-ticks"` // CPUs for adaptive ticks
}

// Param defines a mapping between a YAML field and a kernel parameter.
type Param struct {
	YAMLName    string
	CmdlineName string
	TransformFn func(interface{}) string // Function to transform value to string
}
