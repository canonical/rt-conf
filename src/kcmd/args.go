package kcmd

import (
	"fmt"

	"github.com/canonical/rt-conf/src/helpers"
	"github.com/canonical/rt-conf/src/models"
	"github.com/canonical/rt-conf/src/system"
)

func ProcessKcmdArgs(c *helpers.InternalConfig) error {
	sys, err := system.DetectBootloader()
	if err != nil {
		return fmt.Errorf("failed to detect system: %v", err)
	}
	switch sys {
	case system.Rpi:
		fmt.Println("Raspberry Pi detected")
		models.UpdateRPi(c)
	case system.Grub:
		fmt.Println("GRUB detected")
		err = models.UpdateGrub(c)
		if err != nil {
			return fmt.Errorf("failed to read file: %v", err)
		}
	default:
		return fmt.Errorf("unsupported bootloader: %v", sys)
	}
	return nil
}
