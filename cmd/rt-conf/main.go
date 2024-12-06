package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/canonical/rt-conf/src/execute"
	"github.com/canonical/rt-conf/src/helpers"
	"github.com/canonical/rt-conf/src/models"
	"github.com/canonical/rt-conf/src/system"
)

const (
	cfgFilePath      = "COMMON_CONFIG_PATH"
	ETC_DEFAULT_GRUB = "/etc/default/grub"
)

func getDefaultConfig() string {
	return os.Getenv(cfgFilePath)
}

func main() {

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

	sys, err := system.DetectSystem()
	if err != nil {
		fmt.Printf("Failed to detect system: %v\n", err)
		os.Exit(1)
	}

	conf, err := helpers.LoadConfigFile(*configPath, *grubDefaultPath)
	if err != nil {
		err := fmt.Errorf("failed to load config file: %v", err)
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	if sys == "raspberry" {
		fmt.Println("Raspberry Pi detected")
		execute.RaspberryConclusion(&conf)
	} else {
		err = models.UpdateGrub(&conf)
		if err != nil {
			err := fmt.Errorf("failed to read file: %v", err)
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}
	}

}
