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

type ErrValidation struct {
	err   string
	exist bool
}

var validationErrors = []ErrValidation{
	isolatecpus:    {err: "\n", exist: false},
	enableDynticks: {err: "\n", exist: false},
	adaptiveCPUs:   {err: "\n", exist: false},
}

// TODO: Think in a way to handle if the user wants empty values

// TODO: This function validates and create the error message,
// TODO: maybe it should be split into two functions to separate the logic
func (m *Model) Validation() bool {
	var err error
	var dyntickMode bool

	// TODO: move this logic to outside this function
	// If focusIndex is out of range, return true
	if m.focusIndex < 0 || m.focusIndex >= len(m.inputs) {
		return true
	}

	log.Println("focusIndex on Validation: ", m.focusIndex)
	value := m.inputs[m.focusIndex].Value()

	// TODO: fetch value from YAML file and SetValue()
	// m.inputs[m.focusIndex].SetValue(value)

	if value == "" {
		return false
	}

	switch m.focusIndex {
	case isolatecpus, adaptiveCPUs:
		err = cpu.ValidateList(value, m.iconf.TotalCPUs)
		log.Println("Isolated CPU List: ", value)
		if err != nil {
			validationErrors[m.focusIndex].err = "ERROR: " + err.Error() + "\n"
			validationErrors[m.focusIndex].exist = true
		} else {
			validationErrors[m.focusIndex].err = "\n"
			validationErrors[m.focusIndex].exist = false
		}

	case enableDynticks:
		dyntickMode, err = strconv.ParseBool(value)
		log.Println("Dyntick Mode: ", dyntickMode)
		if err != nil {
			validationErrors[enableDynticks].err =
				"ERROR: expected a boolean value (true|false) got: " +
					value + "\n"
			validationErrors[enableDynticks].exist = true
		} else {
			validationErrors[enableDynticks].err = "\n"
			validationErrors[enableDynticks].exist = false
		}

	}

	// TODO FIX THIS - add the logic to clean up the error message
	m.errorMsg = ""
	for _, v := range validationErrors {
		m.errorMsg += v.err
	}

	for _, v := range validationErrors {
		if v.exist {
			return false
		}
	}

	m.infoMsg = "Press 'enter' to apply changes\n"
	return true
}
