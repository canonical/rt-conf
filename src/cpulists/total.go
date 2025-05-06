package cpulists

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

var execCommand = func(name string, arg ...string) ([]byte, error) {
	return exec.Command(name, arg...).Output()
}

var totalCPUs = func() (int, error) {
	output, err := execCommand("lscpu", "--json")
	if err != nil {
		return 0, err
	}

	var result struct {
		CPUs []struct {
			Field    string `json:"field"`
			Data     string `json:"data"`
			Children []any  `json:"children"`
		} `json:"lscpu"`
	}
	if err := json.Unmarshal(output, &result); err != nil {
		return 0, err
	}

	for _, cpu := range result.CPUs {
		if strings.TrimSpace(cpu.Field) == "CPU(s):" {
			totalCPUs, err := strconv.Atoi(cpu.Data)
			if err != nil {
				return 0, err
			}
			return totalCPUs, nil
		}
	}

	return 0, fmt.Errorf("could not find total CPUs")
}
