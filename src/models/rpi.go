package models

import (
	"github.com/canonical/rt-conf/src/data"
	"github.com/canonical/rt-conf/src/execute"
)

// InjectToGrubFiles inject the kernel command line parameters to the grub files. /etc/default/grub
func UpdateRPi(cfg *data.InternalConfig) ([]string, error) {
	cmdline, err := data.ConstructKeyValuePairs(&cfg.Data.KernelCmdline)
	if err != nil {
		return nil, err
	}
	kcmds := data.DisassembleParamMap(cmdline)
	msgs := execute.RpiConclusion(kcmds)
	return msgs, nil
}
