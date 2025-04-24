package model

import (
	"fmt"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

// TODO: this should be dropped in favor or parsing the default grub file
//
//	** NOTE: instead of using a hardcoded pattern
//	** we should source the /etc/default/grub file, since it's basically
//	** an environment file (key=value) and we can parse it with the
//	** `os` package
var PatternGrubDefault = regexp.MustCompile(`^(GRUB_CMDLINE_LINUX=")([^"]*)(")$`)

// TODO: This need to be changed, and etc/default/grub should be parsed
type FileTransformer interface {
	TransformLine(string) string
	GetFilePath() string
	GetPattern() *regexp.Regexp
}

type Grub struct {
	File    string
	Pattern *regexp.Regexp
}

type Core interface {
	InjectToFile(pattern *regexp.Regexp) error
}

func ReadYAML(path string) (cfg *Config, err error) {
	d, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	err = yaml.Unmarshal([]byte(d), &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %v", err)
	}
	if cfg == nil {
		return nil, fmt.Errorf("empty config file")
	}
	err = cfg.Validate()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
