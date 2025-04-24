package model

import (
	"bufio"
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

// processFile processes a file with a given FileTransformer, applying
// its transformation on lines matching the pattern.
var ProcessFile = func(transformer FileTransformer) error {
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
			// This is where the kcmdline params of bootloader file are updated
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
