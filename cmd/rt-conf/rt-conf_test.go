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
	grubCfgSample = `
# Sample grub.cfg content
	menuentry 'Ubuntu, with Linux 6.8.0-45-generic' --class ubuntu --class gnu-linux --class gnu --class os $menuentry_id_option 'gnulinux-6.8.0-45-generic-advanced-f93ca2dd-74af-4ee8-b478-970b29be5ca3' {
		recordfail
		load_video
		gfxmode $linux_gfx_mode
		insmod gzio
		if [ x$grub_platform = xxen ]; then insmod xzio; insmod lzopio; fi
		insmod part_gpt
		insmod ext2
		search --no-floppy --fs-uuid --set=root f93ca2dd-74af-4ee8-b478-970b29be5ca3
		echo	'Loading Linux 6.8.0-45-generic ...'
		linux	/boot/vmlinuz-6.8.0-45-generic root=UUID=f93ca2dd-74af-4ee8-b478-970b29be5ca3 ro quiet splash $vt_handoff no_hz
		echo	'Loading initial ramdisk ...'
		initrd	/boot/initrd.img-6.8.0-45-generic
	}
`
	grubSample = `
	GRUB_DEFAULT=0
	GRUB_TIMEOUT_STYLE=hidden
	GRUB_TIMEOUT=0
	GRUB_CMDLINE_LINUX_DEFAULT="quiet splash"
	GRUB_CMDLINE_LINUX=""
`

	configSample = `
kernel_cmdline:
  - isolcpus=8-9
  - no_hz
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
	tempGrubCfgPath := setupTempFile(t, grubCfgSample)
	tempGrubPath := setupTempFile(t, grubSample)
	defer os.Remove(tempConfigPath)
	defer os.Remove(tempGrubCfgPath)
	defer os.Remove(tempGrubPath)

	// Prepare InternalConfig with the temporary file paths
	iCfg := helpers.InternalConfig{
		ConfigFile: tempConfigPath,
		GrubCfg: data.Grub{
			File:    tempGrubCfgPath,
			Pattern: regexp.MustCompile(regexGrubcfg),
		},
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

	// Verify grub.cfg updates
	updatedGrubCfg, err := os.ReadFile(tempGrubCfgPath)
	if err != nil {
		t.Fatalf("Failed to read modified grub.cfg: %v", err)
	}
	if !strings.Contains(string(updatedGrubCfg), "isolcpus=8-9") {
		t.Errorf("Expected isolcpus=8-9 in grub.cfg, but not found")
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
