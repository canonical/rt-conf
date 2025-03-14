package models

import (
	"github.com/canonical/rt-conf/src/data"
)

// InjectToGrubFiles inject the kernel command line parameters to the grub files. /etc/default/grub
func UpdateRPi(cfg *data.InternalConfig) ([]string, error) {
	cmdline, err := data.ConstructKeyValuePairs(&cfg.Data.KernelCmdline)
	if err != nil {
		return nil, err
	}
	kcmds := data.ParamsToCmdline(cmdline)
	return RpiConclusion(kcmds), nil
}
