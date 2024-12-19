package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/canonical/rt-conf/src/data"
	"github.com/canonical/rt-conf/src/helpers"
	"github.com/canonical/rt-conf/src/interrupts"
	"github.com/canonical/rt-conf/src/kcmd"
	"github.com/canonical/rt-conf/src/ui"
	tea "github.com/charmbracelet/bubbletea"
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

	runningAsService := flag.Bool("service", false, "Run as a service")

	tui := flag.Bool("ui", false, "Render the TUI")

	flag.Parse()
	if *configPath == "" {
		flag.PrintDefaults()
		log.Fatalf("Failed to load config file: config path not set")
	}

	conf := *data.NewInternCfg()
	if d, err := helpers.LoadConfigFile(*configPath); err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	} else {
		conf.Data = *d
	}

	abs, err := filepath.Abs(*configPath)
	if err != nil {
		log.Fatalf("failed to get absolute path for config file: %v", err)
	}

	conf.CfgFile = abs
	conf.GrubDefault = data.Grub{
		File:    *grubDefaultPath,
		Pattern: data.PatternGrubDefault,
	}

	if *tui {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			log.Fatalf("failed to open log file: %v", err)
		}
		defer f.Close()

		log.Println("Running TUI...")
		log.Println()
		// Run the Terminal User Interface (TUI)
		if _, err := tea.NewProgram(ui.NewModel(&conf), tea.WithAltScreen()).Run(); err != nil {
			log.Fatalf("rt-conf failed: %v", err)
		}
		return
	}

	err = interrupts.ProcessIRQIsolation(&conf)
	if err != nil {
		log.Fatalf("Failed to process interrupts: %v", err)
	}

	// NOTE: This should also be the decision to rather render or not the TUI
	// in the future
	if *runningAsService {
		log.Println("Running as a service")
		return
	}

	// If not running as a service then process the kernel cmdline args
	if err := kcmd.ProcessKcmdArgs(&conf); err != nil {
		log.Fatalf("Failed to process kernel cmdline args: %v", err)
	}

}
