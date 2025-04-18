package kcmd

import (
	"fmt"

	"github.com/canonical/rt-conf/src/model"
)

// InjectToGrubFiles inject the kernel command line parameters to the grub files. /etc/default/grub
func UpdateRPi(cfg *model.InternalConfig) ([]string, error) {
	cmdline := model.ConstructKeyValuePairs(&cfg.Data.KernelCmdline)
	if len(cmdline) == 0 {
		return nil, fmt.Errorf("no parameters to inject")
	}
	kcmds := model.ParamsToCmdline(cmdline)
	return RpiConclusion(kcmds), nil
}
