package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/canonical/rt-conf/src/debug"
	"github.com/canonical/rt-conf/src/irq"
	"github.com/canonical/rt-conf/src/kcmd"
	"github.com/canonical/rt-conf/src/model"
	pwrmgmt "github.com/canonical/rt-conf/src/pwr_mgmt"
)

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	configPath := flags.String("file",
		"",
		"Path to the configuration file")
	grubConfigPath := flags.String("grub-file",
		"/etc/default/grub",
		"Path to the grub configuration file, relevant only for GRUB bootloader")
	verbose := flags.Bool("verbose",
		false,
		"Verbose mode, prints more information to the console")

	if err := flags.Parse(args[1:]); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	if *verbose {
		fmt.Println("Verbose mode enabled")
		debug.Enable()
	}

	if *configPath == "" {
		flag.PrintDefaults()
		err := fmt.Errorf("failed to load config file: path not set")
		log.Println(err)
		return err
	}

	fmt.Println("Configuration file:", *configPath)

	var conf model.InternalConfig
	if d, err := model.LoadConfigFile(*configPath); err != nil {
		err = fmt.Errorf("failed to load config file: %w", err)
		log.Println(err.Error())
		return err
	} else {
		conf.Data = *d
	}

	conf.GrubDefault = model.Grub{
		File: *grubConfigPath,
	}

	// If not running as a service then process the kernel cmdline args
	if msgs, err := kcmd.ProcessKcmdArgs(&conf); err != nil {
		err := fmt.Errorf("failed to process kernel cmdline args: %v", err)
		log.Println(err)
		return err
	} else {
		for _, msg := range msgs {
			fmt.Print(msg)
		}
	}

	err := irq.ApplyIRQConfig(&conf)
	if err != nil {
		err = fmt.Errorf("failed to process interrupts: %v", err)
		log.Println(err)
		return err
	}

	err = pwrmgmt.ApplyPwrConfig(&conf)
	if err != nil {
		err = fmt.Errorf("failed to process power management config: %v", err)
		log.Println(err)
		return err
	}

	return nil
}
