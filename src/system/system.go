package system

import (
	"os"
	"runtime"
	"strings"
)

type System interface {
	DetectHostHardware()
}

const (
	PROC_DT_MODEL = "/proc/device-tree/model"
)

func DetectSystem() (string, error) {
	// Use runtime.GOARCH for a basic architecture check
	if runtime.GOARCH == "arm" || runtime.GOARCH == "arm64" {
		// Check for Raspberry Pi specific file or data
		if _, err := os.Stat(PROC_DT_MODEL); err != nil {
			return runtime.GOARCH, nil
		}

		content, err := os.ReadFile(PROC_DT_MODEL)
		if err != nil {
			return "", err
		}

		if strings.Contains(string(content), "Raspberry Pi") {
			return "raspberry", nil
		}
	}
	return runtime.GOARCH, nil
}
