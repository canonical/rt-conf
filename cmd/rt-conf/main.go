package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/canonical/go-snapctl/env"
	"github.com/canonical/rt-conf/src/debug"
	"github.com/canonical/rt-conf/src/irq"
	"github.com/canonical/rt-conf/src/kcmd"
	"github.com/canonical/rt-conf/src/model"
	pwrmgmt "github.com/canonical/rt-conf/src/pwr_mgmt"
)

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal("Error: ", err)
	}
}

func run(args []string) error {
	envConfigFile := os.Getenv("CONFIG_FILE")
	verboseDefaultCfg := false
	var err error
	envVerbose, ok := os.LookupEnv("VERBOSE")
	if ok {
		verboseDefaultCfg, err = strconv.ParseBool(envVerbose)
		if err != nil {
			return err
		}
	}

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	configPath := flags.String("file",
		envConfigFile,
		"Path to the configuration file")
	grubConfigPath := flags.String("grub-file",
		"/etc/default/grub",
		"Path to the default input grub configuration file, relevant only for GRUB bootloader")
	grubCfgPath := flags.String("grub-custom-file",
		"/etc/default/grub.d/60_rt-conf.cfg",
		"Path to the output drop-in grub configuration file, relevant only for GRUB bootloader")
	verbose := flags.Bool("verbose",
		verboseDefaultCfg,
		"Verbose mode, prints more information to the console")

	if err := flags.Parse(args[1:]); err != nil {
		return fmt.Errorf("failed to parse flags: %v", err)
	}

	log.SetFlags(0)

	if *verbose {
		fmt.Println("Verbose mode enabled")
		debug.Enable()
	}

	if *configPath == "" {
		flag.PrintDefaults()
		return fmt.Errorf("failed to load config file: path not set")
	}

	var conf model.InternalConfig

	if err := conf.Data.LoadFromFile(*configPath); err != nil {
		return fmt.Errorf("failed to load config file: %w", err)
	}

	// If running as a snap, override config with snap options
	if env.Snap() != "" {
		if err := conf.Data.LoadSnapOptions(); err != nil {
			return fmt.Errorf("failed to load config from snap options: %v", err)
		}
	}

	conf.GrubCfg = model.Grub{
		GrubDefaultFilePath: *grubConfigPath,
		CustomGrubFilePath:  *grubCfgPath,
	}

	if msgs, err := kcmd.ProcessKcmdArgs(&conf); err != nil {
		return fmt.Errorf("failed to process kernel cmdline args: %v", err)
	} else {
		for _, msg := range msgs {
			fmt.Print(msg)
		}
	}

	if err := irq.ApplyIRQConfig(&conf); err != nil {
		return fmt.Errorf("failed to process interrupts: %v", err)
	}

	if err := pwrmgmt.ApplyPwrConfig(&conf); err != nil {
		return fmt.Errorf("failed to process power management config: %v", err)
	}

	return nil
}
