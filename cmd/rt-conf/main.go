package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/canonical/rt-conf/src/debug"
	"github.com/canonical/rt-conf/src/irq"
	"github.com/canonical/rt-conf/src/kcmd"
	"github.com/canonical/rt-conf/src/model"
	pwrmgmt "github.com/canonical/rt-conf/src/pwr_mgmt"
)

func main() {
	configPath := flag.String("file",
		"",
		"Path to the configuration file")
	grubConfigPath := flag.String("grub-file",
		"/etc/default/grub",
		"Path to the grub configuration file, relevant only for GRUB bootloader")
	verbose := flag.Bool("verbose",
		false,
		"Verbose mode, prints more information to the console")

	flag.Parse()

	fmt.Println("Reading configuration file from", *configPath)

	if *verbose {
		fmt.Println("Verbose mode enabled")
		debug.Enable()
	}

	if *configPath == "" {
		flag.PrintDefaults()
		log.Fatalf("Failed to load config file: path not set")
	}

	var conf model.InternalConfig
	if d, err := model.LoadConfigFile(*configPath); err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	} else {
		conf.Data = *d
	}

	conf.GrubDefault = model.Grub{
		File: *grubConfigPath,
	}

	// If not running as a service then process the kernel cmdline args
	if msgs, err := kcmd.ProcessKcmdArgs(&conf); err != nil {
		log.Fatalf("Failed to process kernel cmdline args: %v", err)
	} else {
		for _, msg := range msgs {
			fmt.Print(msg)
		}
	}

	err := irq.ApplyIRQConfig(&conf)
	if err != nil {
		log.Fatalf("Failed to process interrupts: %v", err)
	}

	err = pwrmgmt.ApplyPwrConfig(&conf)
	if err != nil {
		log.Fatalf("Failed to process power management config: %v", err)
	}

}
