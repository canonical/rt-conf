package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/canonical/rt-conf/src/data"
	"github.com/canonical/rt-conf/src/helpers"
	"github.com/canonical/rt-conf/src/kcmd"
)

const (
	grubSample = `
GRUB_DEFAULT=0
GRUB_TIMEOUT_STYLE=hidden
GRUB_TIMEOUT=0
GRUB_CMDLINE_LINUX_DEFAULT="quiet splash"
GRUB_CMDLINE_LINUX=""
`

	configSample = `
kernel-cmdline:
  isolcpus: "8-9"
  dyntick-idle: true
  adaptive-ticks: "8-9"
`
)

func setupTempFile(t *testing.T, content string) string {
	tmpFile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	return tmpFile.Name()
}

func TestInjectToFile(t *testing.T) {
	// Set up temporary config-file, grub.cfg and default-grub files
	tempConfigPath := setupTempFile(t, configSample)
	tempGrubPath := setupTempFile(t, grubSample)
	defer os.Remove(tempConfigPath)
	defer os.Remove(tempGrubPath)

	fmt.Printf("tempConfigPath: %s\n", tempConfigPath)
	fmt.Printf("tempGrubPath: %s\n", tempGrubPath)

	var err error
	conf := *data.NewInternCfg()
	if d, err := helpers.LoadConfigFile(tempConfigPath); err != nil {
		t.Fatalf("Failed to load config file: %v", err)
	} else {
		conf.Data = *d
	}

	conf.CfgFile = tempConfigPath
	conf.GrubDefault = data.Grub{
		File:    tempGrubPath,
		Pattern: data.PatternGrubDefault,
	}

	fmt.Printf("Config: %+v\n", conf)

	// Run the InjectToFile method
	_, err = kcmd.ProcessKcmdArgs(&conf) // TODO: Fix this failing step
	if err != nil {
		t.Fatalf("ProcessKcmdArgs failed: %v", err)
	}

	// Verify default-grub updates
	updatedGrub, err := os.ReadFile(tempGrubPath)
	if err != nil {
		t.Fatalf("Failed to read modified grub file: %v", err)
	}

	testCases := []struct {
		param string
		value string
	}{
		{"isolcpus", "8-9"},
		{"nohz", "on"},
		{"nohz_full", "8-9"},
	}
	for _, tc := range testCases {
		if !strings.Contains(string(updatedGrub), fmt.Sprintf("%s=%s", tc.param, tc.value)) {
			t.Errorf("\nExpected %s=%s in grub file, but not found", tc.param, tc.value)
		}
	}
}
