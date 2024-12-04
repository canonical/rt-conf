package models

import "regexp"

type RpiTransformer struct {
	FilePath string
	Pattern  *regexp.Regexp
	Params   []string
}

func (r *RpiTransformer) TransformLine(line string) string {
	// Append each kernel command line parameter to the matched line
	for _, param := range r.Params {
		line += " " + param
	}
	return line
}

func (r *RpiTransformer) GetFilePath() string {
	return r.FilePath
}

func (r *RpiTransformer) GetPattern() *regexp.Regexp {
	return r.Pattern
}
