package models

import (
	"fmt"
	"regexp"
	"strings"
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
	// TODO: Add functionality to avoid duplications of parameters

	// Extract existing parameters
	matches := g.Pattern.FindStringSubmatch(line)
	// Append new parameters
	updatedParams := strings.TrimSpace(matches[2] + " " + strings.Join(g.Params, " "))
	// Reconstruct the line with updated parameters
	return fmt.Sprintf(`%s%s%s`, matches[1], updatedParams, matches[3])
}

func (g *GrubDefaultTransformer) GetFilePath() string {
	return g.FilePath
}

func (g *GrubDefaultTransformer) GetPattern() *regexp.Regexp {
	return g.Pattern
}
