package cpu

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func TotalAvailable() (int, error) {
	cmd := exec.Command("lscpu", "--json")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	var result struct {
		CPUs []struct {
			Field    string        `json:"field"`
			Data     string        `json:"data"`
			Children []interface{} `json:"children"`
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
