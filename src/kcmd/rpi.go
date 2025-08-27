package kcmd

import (
	"fmt"

	"github.com/canonical/rt-conf/src/model"
)

// InjectToGrubFiles inject the kernel command line parameters to the grub files. /etc/default/grub
func UpdateRPi(cfg *model.InternalConfig) ([]string, error) {
	if len(cfg.Data.KernelCmdline.Parameters) == 0 {
		return nil, fmt.Errorf("no parameters to inject")
	}

	return RpiConclusion(cfg.Data.KernelCmdline.Join()), nil
}
