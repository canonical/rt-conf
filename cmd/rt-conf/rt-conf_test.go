package main

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/canonical/rt-conf/src/data"
	"github.com/canonical/rt-conf/src/helpers"
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
	tmpFile, err := os.CreateTemp("", "grub-test-*.cfg")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer tmpFile.Close()

	_, err = tmpFile.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}

	return tmpFile.Name()
}

func TestInjectToFile(t *testing.T) {
	// Set up temporary config-file, grub.cfg and default-grub files
	tempConfigPath := setupTempFile(t, configSample)
	tempGrubPath := setupTempFile(t, grubSample)
	defer os.Remove(tempConfigPath)
	defer os.Remove(tempGrubPath)

	// Prepare InternalConfig with the temporary file paths
	iCfg := helpers.InternalConfig{
		ConfigFile: tempConfigPath,
		GrubDefault: data.Grub{
			File:    tempGrubPath,
			Pattern: regexp.MustCompile(regexGrubDefault),
		},
	}

	// Run the InjectToFile method
	err := iCfg.InjectToGrubFiles()
	if err != nil {
		t.Fatalf("InjectToFile failed: %v", err)
	}

	// Verify default-grub updates
	updatedGrub, err := os.ReadFile(tempGrubPath)
	if err != nil {
		t.Fatalf("Failed to read modified grub file: %v", err)
	}
	if !strings.Contains(string(updatedGrub), "isolcpus=8-9") {
		t.Errorf("Expected isolcpus=8-9 in grub file, but not found")
	}
}
