package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/canonical/rt-conf/src/irq"
	"github.com/canonical/rt-conf/src/kcmd"
	"github.com/canonical/rt-conf/src/model"
	pwrmgmt "github.com/canonical/rt-conf/src/pwr_mgmt"
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

	var conf model.InternalConfig
	if d, err := model.LoadConfigFile(*configPath); err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	} else {
		conf.Data = *d
	}

	abs, err := filepath.Abs(*configPath)
	if err != nil {
		log.Fatalf("failed to get absolute path for config file: %v", err)
	}

	conf.CfgFile = abs
	conf.GrubDefault = model.Grub{
		File:    *grubDefaultPath,
		Pattern: model.PatternGrubDefault,
	}

	// If not running as a service then process the kernel cmdline args
	if msgs, err := kcmd.ProcessKcmdArgs(&conf); err != nil {
		log.Fatalf("Failed to process kernel cmdline args: %v", err)
	} else {
		for _, msg := range msgs {
			fmt.Print(msg)
		}
	}

	err = irq.ApplyIRQConfig(&conf)
	if err != nil {
		log.Fatalf("Failed to process interrupts: %v", err)
	}

	err = pwrmgmt.ApplyPwrConfig(&conf)
	if err != nil {
		log.Fatalf("Failed to process power management config: %v", err)
	}

}
