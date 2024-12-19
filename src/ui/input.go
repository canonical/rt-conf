package ui

import (
	"log"
	"strconv"

	"github.com/canonical/rt-conf/src/cpu"
)

const (
	isolatecpus = iota
	enableDynticks
	adaptiveCPUs
)

func (m *Model) Validation() {
	var err error
	var dyntickMode bool

	switch m.focusIndex {
	case isolatecpus:
		cpuList := m.inputs[m.focusIndex].Value()
		_, err = cpu.ParseCPUs(cpuList, m.iconf.TotalCPUs)

	case enableDynticks:
		dyntickMode, err = strconv.ParseBool(m.inputs[m.focusIndex].Value())
		log.Println("Dyntick Mode: ", dyntickMode)

	case adaptiveCPUs:
		cpuList := m.inputs[m.focusIndex].Value()
		_, err = cpu.ParseCPUs(cpuList, m.iconf.TotalCPUs)
		log.Println("Adaptive ticks CPU List: ", cpuList)

	}

	if err != nil {
		m.errorMsg = "ERROR: " + err.Error()
	}
}
