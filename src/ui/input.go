package ui

import (
	"log"

	"github.com/canonical/rt-conf/src/cpu"
)

// TODO: move this file to a separate module

// NOTE: This function it will only panic during development
func init() {
	// TODO: Find a better way of doing this
	// the struct validationErrors and placeholders_text are coupled
	// they both should have the same length, since it's about the
	// number of kernel parameters to be inserted

	NumParams = len(validationErrorsKcmd)

	if len(validationErrorsKcmd) != NumParams {
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

var validationErrorsKcmd = []ErrValidation{
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
	if m.kcmd.FocusIndex < 0 || m.kcmd.FocusIndex >= len(m.kcmd.Inputs) {
		return validationErrorsKcmd
	}

	log.Println("focusIndex on Validation: ", m.kcmd.FocusIndex)
	value := m.kcmd.Inputs[m.kcmd.FocusIndex].Value()

	// TODO: fetch value from YAML file and SetValue()
	// m.inputs[m.focusIndex].SetValue(value)
	if value == "" {
		validationErrorsKcmd[m.kcmd.FocusIndex].err = "\n"
		validationErrorsKcmd[m.kcmd.FocusIndex].exist = false
	} else {
		m.checkInputs(value)
	}

	m.errorMsg = ""
	for _, v := range validationErrorsKcmd {
		m.errorMsg += v.err
	}

	return validationErrorsKcmd
}

func (m *Model) checkInputs(value string) {
	var err error

	switch m.kcmd.FocusIndex {
	case isolcpusIndex, nohzFullIndex, kthreadsCPUsIndex, irqaffinityIndex:
		err = cpu.ValidateList(value)
		log.Printf("%v: %v ", validationErrorsKcmd[m.kcmd.FocusIndex].name, value)
		if err != nil {
			validationErrorsKcmd[m.kcmd.FocusIndex].err = "ERROR: " + err.Error() + "\n"
			validationErrorsKcmd[m.kcmd.FocusIndex].exist = true
		} else {
			validationErrorsKcmd[m.kcmd.FocusIndex].err = "\n"
			validationErrorsKcmd[m.kcmd.FocusIndex].exist = false
		}

	case nohzIndex:
		if value == "on" {
			validationErrorsKcmd[nohzIndex].err = "\n"
			validationErrorsKcmd[nohzIndex].exist = false
		} else if value == "off" {
			validationErrorsKcmd[nohzIndex].err = "\n"
			validationErrorsKcmd[nohzIndex].exist = false
		} else {
			validationErrorsKcmd[nohzIndex].err =
				"ERROR: expected (on) or (off) value got: " + value + "\n"
			validationErrorsKcmd[nohzIndex].exist = true
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
