package main

import (
	"flag"
	"log"
	"os"

	"github.com/canonical/rt-conf/src/helpers"
	"github.com/canonical/rt-conf/src/kcmd"
)

const (
	cfgFilePath      = "COMMON_CONFIG_PATH"
	ETC_DEFAULT_GRUB = "/etc/default/grub"
)

const (
	cfgFilehelp = `Path to the configuration file either set this or set the COMMON_CONFIG_PATH environment variable`
	grubHelp    = `Path to the default grub file`
)

func getDefaultConfig() string {
	return os.Getenv(cfgFilePath)
}

func main() {
	configPath := flag.String("config", getDefaultConfig(), cfgFilehelp)

	// TODO: make this generic for any bootloader
	// Define the paths to grub as flags
	grubDefaultPath := flag.String("grub-default", ETC_DEFAULT_GRUB, grubHelp)

	flag.Parse()
	if *configPath == "" {
		flag.PrintDefaults()
		log.Fatalf("Failed to load config file: config path not set")
	}

	conf, err := helpers.LoadConfigFile(*configPath, *grubDefaultPath)
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}

	err = kcmd.ProcessKcmdArgs(&conf)
	if err != nil {
		log.Fatalf("Failed to process kernel cmdline args: %v", err)
	}
}
