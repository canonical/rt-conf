package main

import "regexp"

type Core interface {
	InjectToFile(pattern *regexp.Regexp) error
}

type grub struct {
	file    string
	pattern *regexp.Regexp
}

type InternalConfig struct {
	configFile string
	data       Config

	grubDefault grub
	grubCfg     grub
}

type Config struct {
	KernelCmdline []string `yaml:"kernel_cmdline"`
}
