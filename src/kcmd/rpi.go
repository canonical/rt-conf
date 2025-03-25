package kcmd

import (
	"github.com/canonical/rt-conf/src/model"
)

// InjectToGrubFiles inject the kernel command line parameters to the grub files. /etc/default/grub
func UpdateRPi(cfg *model.InternalConfig) ([]string, error) {
	cmdline, err := model.ConstructKeyValuePairs(&cfg.Data.KernelCmdline)
	if err != nil {
		return nil, err
	}
	kcmds := model.ParamsToCmdline(cmdline)
	return RpiConclusion(kcmds), nil
}
