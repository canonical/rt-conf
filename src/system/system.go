package system

import (
	"os"
	"strings"
)

type SystemType int

const (
	Unknown SystemType = iota
	Grub
	Rpi
	Uboot
	UbuntuCore
)

var baseDir = "" // baseDir is used to mock the file system in tests

var DetectSystem = func() (SystemType, error) {
	// Verify if SNAP_SAVE_DATA is present, indicating the system is Ubuntu Core
	// see: https://snapcraft.io/docs/environment-variables#heading--snap-save-data
	_, isUC := os.LookupEnv("SNAP_SAVE_DATA")
	if isUC {
		return UbuntuCore, nil
	}

	if _, err := os.Stat(baseDir + "/proc/device-tree/model"); err == nil {
		content, err := os.ReadFile("/proc/device-tree/model")
		if err != nil {
			return Unknown, err
		}

		if strings.Contains(string(content), "Raspberry Pi") {
			return Rpi, nil
		}
		return Unknown, nil
	}

	if _, err := os.Stat(baseDir + "/etc/default/grub"); err == nil {
		return Grub, nil
	}

	return Unknown, nil
}
