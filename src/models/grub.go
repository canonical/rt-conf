package models

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/canonical/rt-conf/src/data"
	"github.com/canonical/rt-conf/src/execute"
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
	Params   []string
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

	// Append the parameters
	updatedParams := strings.TrimSpace(" " + strings.Join(g.Params, " "))

	// Reconstruct the line with updated parameters
	return fmt.Sprintf(`%s%s%s`, matches[1], updatedParams, matches[3])
}

func (g *GrubDefaultTransformer) GetFilePath() string {
	return g.FilePath
}

func (g *GrubDefaultTransformer) GetPattern() *regexp.Regexp {
	return g.Pattern
}

// InjectToGrubFiles inject the kernel command line parameters to the grub files. /etc/default/grub
func UpdateGrub(cfg *data.InternalConfig) ([]string, error) {
	var msgs []string

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
	paramsStrLine, ok := grubMap["GRUB_CMDLINE_LINUX"]
	if !ok {
		return nil,
			fmt.Errorf("GRUB_CMDLINE_LINUX not found in %s",
				grubDefault.FilePath)
	}
	currParams := data.AssembleParamMap(paramsStrLine)

	// This replaces if the param already exists and
	// creates a new one if it doesn't
	for k, v := range params {
		currParams[k] = v
	}
	grubDefault.Params = data.DisassembleParamMap(currParams)
	log.Println("Final kcmdline: ", grubDefault.Params)

	if err := data.ProcessFile(grubDefault); err != nil {
		return nil, fmt.Errorf("error updating %s: %v", grubDefault.FilePath, err)
	}

	msgs = append(msgs, "Updated default grub file: "+grubDefault.FilePath+"\n")
	msgs = append(msgs, execute.GrubConclusion()...)

	return msgs, nil
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
