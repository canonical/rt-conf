package kcmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/canonical/rt-conf/src/model"
)

// grubDefaultTransformer handles transformations for /etc/default/grub
type GrubDefaultTransformer struct {
	FilePath string
	Pattern  *regexp.Regexp
	Cmdline  string
}

func (g *GrubDefaultTransformer) TransformLine(line string) string {
	// Extract existing parameters
	matches := g.Pattern.FindStringSubmatch(line)

	// Reconstruct the line with updated parameters
	return fmt.Sprintf(`%s%s%s`, matches[1], g.Cmdline, matches[3])
}

func (g *GrubDefaultTransformer) GetFilePath() string {
	return g.FilePath
}

func (g *GrubDefaultTransformer) GetPattern() *regexp.Regexp {
	return g.Pattern
}

// InjectToGrubFiles inject the kernel command line parameters to the grub files. /etc/default/grub
func UpdateGrub(cfg *model.InternalConfig) ([]string, error) {

	params := model.ConstructKeyValuePairs(&cfg.Data.KernelCmdline)
	if len(params) == 0 {
		return nil, fmt.Errorf("no parameters to inject")
	}
	grubDefault := &GrubDefaultTransformer{
		FilePath: cfg.GrubDefault.File,
		Pattern:  model.PatternGrubDefault,
	}

	grubMap, err := ParseDefaultGrubFile(grubDefault.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse grub file: %v", err)
	}
	cmdline, ok := grubMap["GRUB_CMDLINE_LINUX"]
	if !ok {
		return nil,
			fmt.Errorf("GRUB_CMDLINE_LINUX not found in %s",
				grubDefault.FilePath)
	}

	if err := duplicatedParams(cmdline); err != nil {
		return nil, fmt.Errorf(
			"invalid existing parameters in %s for GRUB_CMDLINE_LINUX: %s",
			grubDefault.FilePath, err)
	}
	currParams := model.CmdlineToParams(cmdline)

	// This replaces if the param already exists and
	// creates a new one if it doesn't
	for k, v := range params {
		currParams[k] = v
	}
	grubDefault.Cmdline = model.ParamsToCmdline(currParams)
	log.Println("Final kcmdline:", grubDefault.Cmdline)

	if err := processFile(grubDefault); err != nil {
		return nil, fmt.Errorf("error updating %s: %v", grubDefault.FilePath, err)
	}

	return GrubConclusion(grubDefault.FilePath), nil
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
var processFile = func(transformer model.FileTransformer) error {
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
