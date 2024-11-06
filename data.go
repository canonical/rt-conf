package main

import "regexp"

type Core interface {
	InjectToFile(pattern *regexp.Regexp) error
}

type InternalConfig struct {
	grubFile           string
	configFile         string
	data               Config
	patternCfgGrub     *regexp.Regexp
	patternDefaultGrug *regexp.Regexp
}

type Config struct {
	KernelCmdline []string `yaml:"kernel_cmdline"`
}
