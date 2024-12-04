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

// Holds the parsed YAML data
type Config map[string]interface{}

// Param defines a mapping between a YAML field and a kernel parameter.
type Param struct {
	YAMLName    string
	CmdlineName string
	TransformFn func(interface{}) string // Function to transform value to string
}
