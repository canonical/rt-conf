package models

import (
	"github.com/canonical/rt-conf/src/execute"
	"github.com/canonical/rt-conf/src/helpers"
)

// InjectToGrubFiles inject the kernel command line parameters to the grub files. /etc/default/grub
func UpdateRPi(cfg *helpers.InternalConfig) {
	cmdline := helpers.TranslateConfig(cfg.Data)
	execute.RpiConclusion(cmdline)
}
