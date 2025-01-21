package ui

import (
	"log"

	"github.com/canonical/rt-conf/src/cpu"
)

// TODO: move this file to a separate module

func init() {
	NumParams = len(validationErrors)

	if len(validationErrors) != NumParams {
		log.Fatalf("Number of validationErrors is different from NumParams")
	}

	if len(placeholders_text) != NumParams {
		log.Fatalf("Number of placeholders_text is different from NumParams")
	}
}

var NumParams int

type ErrValidation struct {
	err   string
	exist bool
	name  string
}

var validationErrors = []ErrValidation{
	isolcpusIndex:     {err: "\n", exist: true, name: "isolcpus"},
	nohzIndex:         {err: "\n", exist: true, name: "nohz"},
	nohzFullIndex:     {err: "\n", exist: true, name: "nohz_full"},
	kthreadsCPUsIndex: {err: "\n", exist: true, name: "kthreadsCPUs"},
	irqaffinityIndex:  {err: "\n", exist: true, name: "irqaffinity"},
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

	return validationErrors
}

func (m *Model) checkInputs(value string) {
	var err error

	switch m.focusIndex {
	case isolcpusIndex, nohzFullIndex, kthreadsCPUsIndex, irqaffinityIndex:
		err = cpu.ValidateList(value)
		log.Printf("%v: %v ", validationErrors[m.focusIndex].name, value)
		if err != nil {
			validationErrors[m.focusIndex].err = "ERROR: " + err.Error() + "\n"
			validationErrors[m.focusIndex].exist = true
		} else {
			validationErrors[m.focusIndex].err = "\n"
			validationErrors[m.focusIndex].exist = false
		}

	case nohzIndex:
		if value == "on" {
			validationErrors[nohzIndex].err = "\n"
			validationErrors[nohzIndex].exist = false
		} else if value == "off" {
			validationErrors[nohzIndex].err = "\n"
			validationErrors[nohzIndex].exist = false
		} else {
			validationErrors[nohzIndex].err =
				"ERROR: expected (on) or (off) value got: " + value + "\n"
			validationErrors[nohzIndex].exist = true
			break
		}

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
