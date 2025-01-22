package data

import (
	"regexp"
)

var PatternGrubDefault = regexp.MustCompile(`^(GRUB_CMDLINE_LINUX=")([^"]*)(")$`)

type InternalConfig struct {
	CfgFile string
	Data    Config

	GrubDefault Grub
}

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
	Interrupts    IRQs          `yaml:"interrupts"`
	KernelCmdline KernelCmdline `yaml:"kernel-cmdline"`
}

type IRQs struct {
	IsolateCPU string `yaml:"remove-from-cpus"`
	IRQHandler string `yaml:"handle-on-cpus"`
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
