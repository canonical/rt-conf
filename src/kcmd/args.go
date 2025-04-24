package kcmd

import (
	"fmt"
	"log"

	"github.com/canonical/rt-conf/src/model"
	"github.com/canonical/rt-conf/src/system"
)

var kcmdSys = map[system.SystemType]func(*model.InternalConfig) ([]string, error){
	system.Rpi:        UpdateRPi,
	system.Grub:       UpdateGrub,
	system.UbuntuCore: UpdateUbuntuCore,
}

func ProcessKcmdArgs(c *model.InternalConfig) ([]string, error) {
	if c.Data.KernelCmdline == (model.KernelCmdline{}) {
		// No kernel command line options to process
		log.Println("No kernel command line options to process")
		return nil, nil
	}

	var msgs []string
	sys, err := system.DetectSystem()
	if err != nil {
		return nil, fmt.Errorf("failed to detect system: %v", err)
	}
	processKcmd, ok := kcmdSys[sys]
	if !ok {
		return nil, fmt.Errorf("unsupported bootloader: %v", sys)
	}
	tmp, err := processKcmd(c)
	if err != nil {
		return nil, err
	}
	msgs = append(msgs, tmp...)

	return msgs, nil
}
