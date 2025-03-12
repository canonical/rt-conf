package system

import (
	"os"
	"strings"
)

func DetectSystem() (SystemType, error) {
	// Verify if SNAP_SAVE_DATA is present, indicating the system is Ubuntu Core
	// see: https://snapcraft.io/docs/environment-variables#heading--snap-save-data
	_, isUC := os.LookupEnv("SNAP_SAVE_DATA")
	if isUC {
		return UCore, nil
	}

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
