package ui

import (
	"log"

	"github.com/canonical/rt-conf/src/cpu"
)

type ErrValidation struct {
	err   string
	exist bool
}

var validationErrors = []ErrValidation{
	isolcpusIndex: {err: "\n", exist: false},
	nohzIndex:     {err: "\n", exist: false},
	nohzFullIndex: {err: "\n", exist: false},
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
		return true
	}

	switch m.focusIndex {
	case isolcpusIndex, nohzFullIndex:
		err = cpu.ValidateList(value, m.iconf.TotalCPUs)
		log.Println("Isolated CPU List: ", value)
		if err != nil {
			validationErrors[m.focusIndex].err = "ERROR: " + err.Error() + "\n"
			validationErrors[m.focusIndex].exist = true
		} else {
			validationErrors[m.focusIndex].err = "\n"
			validationErrors[m.focusIndex].exist = false
		}

	case nohzIndex:

		if value == "y" || value == "Y" {
			value = "true"
		} else if value == "n" || value == "N" {
			value = "false"
		} else {
			validationErrors[nohzIndex].err =
				"ERROR: expected Yes or No value (y|n) got: " + value + "\n"
			validationErrors[nohzIndex].exist = true
			break
		}

		log.Println("Dyntick Mode: ", dyntickMode)
		if err != nil {
			validationErrors[nohzIndex].err =
				"ERROR: expected Yes or No value (y|n) got: " +
					value + "\n"
			validationErrors[nohzIndex].exist = true
		} else {
			validationErrors[nohzIndex].err = "\n"
			validationErrors[nohzIndex].exist = false
		}
	default:
		break

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
