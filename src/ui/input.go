package ui

import (
	"fmt"
	"log"
	"strings"

	"github.com/canonical/rt-conf/src/cpulist"
	"github.com/canonical/rt-conf/src/data"
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
func (m *KcmdlineMenuModel) Validation() []ErrValidation {

	// TODO: move this logic to outside this function
	// If focusIndex is out of range, just return the validationErrors
	if m.FocusIndex < 0 || m.FocusIndex >= len(m.Inputs) {
		return validationErrorsKcmd
	}

	log.Println("focusIndex on Validation: ", m.FocusIndex)
	value := m.Inputs[m.FocusIndex].Value()

	// TODO: fetch value from YAML file and SetValue()
	// m.inputs[m.focusIndex].SetValue(value)
	if value == "" {
		validationErrorsKcmd[m.FocusIndex].err = "\n"
		validationErrorsKcmd[m.FocusIndex].exist = false
	} else {
		m.checkInputs(value)
	}

	m.errorMsg = ""
	for _, v := range validationErrorsKcmd {
		m.errorMsg += v.err
	}

	return validationErrorsKcmd
}

func (m *KcmdlineMenuModel) checkInputs(value string) {
	var err error

	switch m.FocusIndex {
	case isolcpusIndex, nohzFullIndex, kthreadsCPUsIndex, irqaffinityIndex:
		err = cpulist.ValidateList(value)
		log.Printf("%v: %v ", validationErrorsKcmd[m.FocusIndex].name, value)
		if err != nil {
			validationErrorsKcmd[m.FocusIndex].err = "ERROR: " + err.Error() + "\n"
			validationErrorsKcmd[m.FocusIndex].exist = true
		} else {
			validationErrorsKcmd[m.FocusIndex].err = "\n"
			validationErrorsKcmd[m.FocusIndex].exist = false
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

func (m *KcmdlineMenuModel) AreValidInputs() bool {
	validated := m.Validation()
	for _, v := range validated {
		if v.exist {
			return false
		}
	}
	return true
}

// ***** IRQ ADD/EDIT VIEW ******//

var validationErrorsIRQ = []ErrValidation{
	irqFilterIndex: {err: "\n", exist: true, name: "irq_filter"},
	cpuListIndex:   {err: "\n", exist: true, name: "cpulist"},
}

// TODO: improve this validt: it's improperlly reporting as valid empty values
func (m *IRQAddEditMenu) Validation() []ErrValidation {
	log.Println("---- IRQEditMode VALIDATION ----")

	// TODO: move this logic to outside this function
	// If focusIndex is out of range, just return the validationErrors

	if m.FocusIndex < 0 || m.FocusIndex >= len(m.Inputs) {
		log.Println("[WARN] FocusIndex out of range")
		return validationErrorsIRQ
	}

	value := m.Inputs[m.FocusIndex].Value()
	log.Println("focusIndex on Validation: ", m.FocusIndex)
	log.Println("value on Validation: ", value)

	if value == "" {
		log.Println("[NOTE] Value is empty")
		validationErrorsIRQ[m.FocusIndex].err = "\n"
		// If a input is empty, it cannot be a valid IRQ rule
		validationErrorsIRQ[m.FocusIndex].exist = true
	} else {
		log.Println("Validating value: ", value)
		m.checkInputs(value)
	}

	m.errorMsg = ""
	for _, v := range validationErrorsIRQ {
		m.errorMsg += v.err
	}

	return validationErrorsIRQ
}

func (m *IRQAddEditMenu) checkInputs(value string) {
	var err error
	log.Println(">>>>>>>>>> checkInputs: ", value)

	switch m.FocusIndex {
	case cpuListIndex:
		err = cpulist.ValidateList(value)
		// log.Println("Validating cpulist: ", value, "err: ", err)
		// log.Printf("%v: %v ", validationErrorsIRQ[m.FocusIndex].name, value)
		if err != nil {
			validationErrorsIRQ[m.FocusIndex].err = "ERROR: " + err.Error() + "\n"
			validationErrorsIRQ[m.FocusIndex].exist = true
		} else {
			validationErrorsIRQ[m.FocusIndex].err = "\n"
			validationErrorsIRQ[m.FocusIndex].exist = false
		}

	case irqFilterIndex:
		// log.Println("Validating irqFilter: ", value)
		irqFilter, err := ParseIRQFilter(value)
		if err != nil {
			validationErrorsIRQ[m.FocusIndex].err = "ERROR: " + err.Error() + "\n"
			validationErrorsIRQ[m.FocusIndex].exist = true
			break
		}
		err = irqFilter.Validate()
		if err != nil {
			validationErrorsIRQ[m.FocusIndex].err = "ERROR: " + err.Error() + "\n"
			validationErrorsIRQ[m.FocusIndex].exist = true
			break
		}

		validationErrorsIRQ[m.FocusIndex].err = "\n"
		validationErrorsIRQ[m.FocusIndex].exist = false
		// TODO: the irqFilter must be transmited

	default:
		log.Println("[ERROR] FocusIndex out of range")
	}
}

func (m *IRQAddEditMenu) AreValidInputs() bool {
	validated := m.Validation()
	for _, v := range validated {
		if v.exist {
			return false
		}
	}
	return true
}

func (m *IRQAddEditMenu) RunInputValidation() {
	m.AreValidInputs()
}

// ParseIRQFilter parses a string like:
// "  actions:<string>  chip_name:<string> name:<string> type:<string>   "
// into an IRQFilter. At least one filter must be provided.
func ParseIRQFilter(query string) (data.IRQFilter, error) {
	var filter data.IRQFilter

	// Trim leading/trailing whitespace
	query = strings.TrimSpace(query)
	if query == "" {
		return filter, fmt.Errorf("query is empty")
	}

	// Split the query string by whitespace.
	// strings.Fields splits on any sequence of whitespace characters.
	tokens := strings.Fields(query)
	for _, token := range tokens {
		// Each token is expected to be in the form "key:value"
		parts := strings.SplitN(token, ":", 2)
		if len(parts) != 2 {
			return filter,
				fmt.Errorf(
					"invalid token: %q, expected format is key:value",
					token,
				)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "actions":
			filter.Actions = value
		case "chip_name":
			filter.ChipName = value
		case "name":
			filter.Name = value
		case "type":
			filter.Type = value
		default:
			return filter,
				fmt.Errorf(
					"unknown key: %q, valid keys are: actions, chip_name, name, type",
					key,
				)
		}
	}

	// Ensure at least one filter field is non-empty.
	if filter.Actions == "" &&
		filter.ChipName == "" &&
		filter.Name == "" &&
		filter.Type == "" {
		return filter, fmt.Errorf("at least one filter field must be provided")
	}
	return filter, nil
}
