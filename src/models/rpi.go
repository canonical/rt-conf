package models

import (
	"github.com/canonical/rt-conf/src/data"
	"github.com/canonical/rt-conf/src/execute"
	"github.com/canonical/rt-conf/src/helpers"
)

// InjectToGrubFiles inject the kernel command line parameters to the grub files. /etc/default/grub
func UpdateRPi(cfg *data.InternalConfig) ([]string, error) {

	cmdline, err := helpers.ReconstructKeyValuePairs(&cfg.Data.KernelCmdline)
	if err != nil {
		return nil, err
	}
	msgs := execute.RpiConclusion(cmdline)
	return msgs, nil
}
