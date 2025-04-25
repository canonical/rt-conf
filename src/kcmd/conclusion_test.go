package kcmd

import "testing"

func TestGrubConclusion(t *testing.T) {
	grubFile := "/etc/default/grub"
	expected := []string{
		"Detected bootloader: GRUB\n",
		"Updated default grub file: " + grubFile + "\n",
		"\n",
		"Please run:\n",
		"\n",
		"\tsudo update-grub\n",
		"\n",
		"to apply the changes to your bootloader.\n",
	}
	result := GrubConclusion(grubFile)

	if len(result) != len(expected) {
		t.Errorf("Expected %d lines, got %d", len(expected), len(result))
	}

	for i, line := range expected {
		if result[i] != line {
			t.Errorf("Expected line %d to be '%s', got '%s'", i, line, result[i])
		}
	}
}

func TestRpiConclusion(t *testing.T) {
	cmdline := "quiet splash isolcpus=1-2"
	expected := []string{
		"Detected bootloader: Raspberry Pi\n",
		"\n",
		"Please, append the following to /boot/firmware/cmdline.txt:\n",
		"In case of old style boot partition,\n",
		"append to /boot/cmdline.txt\n",
		cmdline,
		"\n",
	}
	result := RpiConclusion(cmdline)

	if len(result) != len(expected) {
		t.Errorf("Expected %d lines, got %d", len(expected), len(result))
	}

	for i, line := range expected {
		if result[i] != line {
			t.Errorf("Expected line %d to be '%s', got '%s'", i, line, result[i])
		}
	}
}

func TestUbuntuCoreConclusion(t *testing.T) {
	expected := []string{
		"Detected bootloader: Ubuntu Core managed\n",
		"\n",
		"Sucessfully applied the changes.\n",
		"Please reboot your system to apply the changes.\n",
	}
	result := UbuntuCoreConclusion()

	if len(result) != len(expected) {
		t.Errorf("Expected %d lines, got %d", len(expected), len(result))
	}

	for i, line := range expected {
		if result[i] != line {
			t.Errorf("Expected line %d to be '%s', got '%s'", i, line, result[i])
		}
	}
}
