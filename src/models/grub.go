package models

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/canonical/rt-conf/src/data"
)

// grubCfgTransformer handles transformations for /boot/grub/grub.cfg
type GrubCfgTransformer struct {
	FilePath string
	Pattern  *regexp.Regexp
	Params   []string
}

// grubDefaultTransformer handles transformations for /etc/default/grub
type GrubDefaultTransformer struct {
	FilePath string
	Pattern  *regexp.Regexp
	Cmdline  string
}

func (g *GrubCfgTransformer) TransformLine(line string) string {
	// Append each kernel command line parameter to the matched line
	for _, param := range g.Params {
		line += " " + param
	}
	return line
}

func (g *GrubCfgTransformer) GetFilePath() string {
	return g.FilePath
}

func (g *GrubCfgTransformer) GetPattern() *regexp.Regexp {
	return g.Pattern
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
func UpdateGrub(cfg *data.InternalConfig) ([]string, error) {

	params, err := data.ConstructKeyValuePairs(&cfg.Data.KernelCmdline)
	if err != nil {
		return nil, fmt.Errorf("failed to reconstruct key-value pairs: %v", err)
	}
	grubDefault := &GrubDefaultTransformer{
		FilePath: cfg.GrubDefault.File,
		Pattern:  cfg.GrubDefault.Pattern,
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
	currParams := data.CmdlineToParams(cmdline)

	// This replaces if the param already exists and
	// creates a new one if it doesn't
	for k, v := range params {
		currParams[k] = v
	}
	grubDefault.Cmdline = data.ParamsToCmdline(currParams)
	log.Println("Final kcmdline:", grubDefault.Cmdline)

	if err := data.ProcessFile(grubDefault); err != nil {
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
		param, ok := params[pair[0]]
		if ok {
			// Skip parameters without a value, they can be safelly dropped
			if len(pair) != 2 {
				// Value is optional for some kernel cmdline parameters
				params[p] = ""
				continue
			}

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
