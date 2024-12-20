package kcmd

import (
	"fmt"

	"github.com/canonical/rt-conf/src/data"
	"github.com/canonical/rt-conf/src/models"
	"github.com/canonical/rt-conf/src/system"
)

func ProcessKcmdArgs(c *data.InternalConfig) ([]string, error) {
	var msgs []string
	sys, err := system.DetectBootloader()
	if err != nil {
		return nil, fmt.Errorf("failed to detect system: %v", err)
	}
	switch sys {
	case system.Rpi:
		msgs = append(msgs, "Raspberry Pi detected\n")
		tmp := models.UpdateRPi(c)
		msgs = append(msgs, tmp...)
	case system.Grub:
		msgs = append(msgs, "GRUB detected\n")
		tmp, err := models.UpdateGrub(c)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %v", err)
		}
		msgs = append(msgs, tmp...)
	default:
		return nil, fmt.Errorf("unsupported bootloader: %v", sys)
	}
	return msgs, nil
}
