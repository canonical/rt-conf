package system

import (
	"os"
	"strings"
)

type System interface {
	DetectHostHardware()
}

// TODO: Return an value from an enum

func DetectBootloader() (Bootloader, error) {
	// Use runtime.GOARCH for a basic architecture check

	// TODO: check for the vendor
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
