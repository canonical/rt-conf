package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"sync"

	"github.com/canonical/rt-conf/src/data"
	"github.com/canonical/rt-conf/src/helpers"
)

const (
	cfgFilePath = "COMMON_CONFIG_PATH"
)

var lock = &sync.Mutex{}

type singletonDefaultConfig struct {
	path string
}

var instance *singletonDefaultConfig

func getDefaultConfig() string {
	lock.Lock()
	defer lock.Unlock()

	return instance.path
}

func init() {
	defaultConfig := ""
	defaultConfig = os.Getenv(cfgFilePath)
	instance = &singletonDefaultConfig{path: defaultConfig}
}

const (
	// Grub files paths
	// BOOT_GRUB_GRUBCFG = "/boot/grub/grub.cfg"
	ETC_DEFAULT_GRUB = "/etc/default/grub"

	// // Default configuration file path
	// DEFAULT_CONFIG_PATH = "/var/snap/rt-conf/common/config.yaml"

	RegexGrubDefault = `^(GRUB_CMDLINE_LINUX=")([^"]*)(")$`

	// regexGrubcfg = `linux\s*\/*\w*\/vmlinuz-\d.\d.\d`
)

func main() {
	// TODO: Add system detection functionality

	configPath := flag.String("config", getDefaultConfig(), "Path to the configuration file")

	// Define the paths to the grub files as flags
	// To be used for testing purposes
	// grubCfgPath := flag.String("grub-cfg", BOOT_GRUB_GRUBCFG, "Path to the processed grub file")
	grubDefaultPath := flag.String("grub-default", ETC_DEFAULT_GRUB, "Path to the default grub file")

	flag.Parse()
	if *configPath == "" {
		fmt.Fprintln(os.Stderr, "Default config path not set neither by flag nor by env var")
		fmt.Fprintf(os.Stderr, "Please set the %v environment variable\n", cfgFilePath)
		fmt.Fprintf(os.Stderr, " or use the --config flag to set the path to the configuration file\n")
		os.Exit(1)
	}

	iCfg := helpers.InternalConfig{
		ConfigFile: *configPath,
		// GrubCfg: data.Grub{
		// 	File:    *grubCfgPath,
		// 	Pattern: regexp.MustCompile(regexGrubcfg),
		// },
		GrubDefault: data.Grub{
			File:    *grubDefaultPath,
			Pattern: regexp.MustCompile(RegexGrubDefault),
		},
	}

	fmt.Println("Config path: ", iCfg.ConfigFile)
	// fmt.Println("Grub path: ", iCfg.GrubCfg.File)

	err := iCfg.InjectToGrubFiles()
	if err != nil {
		fmt.Printf("Failed to inject to file: %v\n", err)
		os.Exit(1)
	}

	// TODO: Add system detection functionality to print the message for each system
	fmt.Println("Successfully injected to file")
	fmt.Println("Please run:\nsudo update-grub\nto apply the changes")
}
