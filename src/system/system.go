package system

import (
	"os"
	"strings"
)

func DetectBootloader() (Bootloader, error) {
	if _, err := os.Stat("/proc/device-tree/model"); err == nil {
		content, err := os.ReadFile("/proc/device-tree/model")
		if err != nil {
			return Unknown, err
		}

		if strings.Contains(string(content), "Raspberry Pi") {
			return Rpi, nil
		}
		return Unknown, nil
	}

	if _, err := os.Stat("/etc/default/grub"); err == nil {
		return Grub, nil
	}

	return Unknown, nil
}
