package ui

import (
	"log"

	"github.com/canonical/rt-conf/src/cpu"
)

// TODO: move this file to a separate module

type ErrValidation struct {
	err   string
	exist bool
}

var validationErrors = []ErrValidation{
	isolcpusIndex: {err: "\n", exist: true},
	nohzIndex:     {err: "\n", exist: true},
	nohzFullIndex: {err: "\n", exist: true},
}

// TODO: Think in a way to handle if the user wants empty values

// TODO: This function validates and create the error message,
func (m *Model) Validation() []ErrValidation {

	// TODO: move this logic to outside this function
	// If focusIndex is out of range, just return the validationErrors
	if m.focusIndex < 0 || m.focusIndex >= len(m.inputs) {
		return validationErrors
	}

	log.Println("focusIndex on Validation: ", m.focusIndex)
	value := m.inputs[m.focusIndex].Value()

	// TODO: fetch value from YAML file and SetValue()
	// m.inputs[m.focusIndex].SetValue(value)
	if value == "" {
		validationErrors[m.focusIndex].err = "\n"
		validationErrors[m.focusIndex].exist = false
	} else {
		m.checkInputs(value)
	}

	m.errorMsg = ""
	for _, v := range validationErrors {
		m.errorMsg += v.err
	}

	m.infoMsg = "Press 'enter' to apply changes\n"
	return validationErrors
}

func (m *Model) checkInputs(value string) {
	var err error
	var dyntickMode bool

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
			validationErrors[nohzIndex].err = "\n"
			validationErrors[nohzIndex].exist = false
			dyntickMode = true
		} else if value == "n" || value == "N" {
			validationErrors[nohzIndex].err = "\n"
			validationErrors[nohzIndex].exist = false
			dyntickMode = false
		} else {
			validationErrors[nohzIndex].err =
				"ERROR: expected Yes or No value (y|n) got: " + value + "\n"
			validationErrors[nohzIndex].exist = true
			break
		}

		log.Println("Dyntick Mode: ", dyntickMode)
	default:
		break
	}
}

func (m *Model) AreValidInputs() bool {
	validated := m.Validation()
	for _, v := range validated {
		if v.exist {
			return false
		}
	}
	return true
}
