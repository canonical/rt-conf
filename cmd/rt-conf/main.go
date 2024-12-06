package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/canonical/rt-conf/src/helpers"
)

const (
	cfgFilePath      = "COMMON_CONFIG_PATH"
	ETC_DEFAULT_GRUB = "/etc/default/grub"
)

func getDefaultConfig() string {
	return os.Getenv(cfgFilePath)
}

func main() {
	// TODO: Add system detection functionality

	configPath := flag.String("config", getDefaultConfig(), "Path to the configuration file")

	// Define the paths to grub as flags
	grubDefaultPath := flag.String("grub-default", ETC_DEFAULT_GRUB, "Path to the default grub file")

	flag.Parse()
	if *configPath == "" {
		fmt.Fprintln(os.Stderr, "Default config path not set neither by flag nor by env var")
		fmt.Fprintf(os.Stderr, "Please set the %v environment variable\n", cfgFilePath)
		fmt.Fprintf(os.Stderr, " or use the --config flag to set the path to the configuration file\n")
		os.Exit(1)
	}

	conf, err := helpers.LoadConfigFile(*configPath, *grubDefaultPath)
	if err != nil {
		fmt.Printf("Failed to load config file: %v\n", err)
		os.Exit(1)
	}

	err = helpers.UpdateGrub(&conf)
	if err != nil {
		fmt.Printf("Failed to inject to file: %v\n", err)
		os.Exit(1)
	}

}
