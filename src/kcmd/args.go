package kcmd

import (
	"fmt"

	"github.com/canonical/rt-conf/src/data"
	"github.com/canonical/rt-conf/src/models"
	"github.com/canonical/rt-conf/src/system"
)

var kcmdSys = map[system.SystemType]func(*data.InternalConfig) ([]string, error){
	system.Rpi:   models.UpdateRPi,
	system.Grub:  models.UpdateGrub,
	system.UCore: models.UpdateCore,
}

func ProcessKcmdArgs(c *data.InternalConfig) ([]string, error) {
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
