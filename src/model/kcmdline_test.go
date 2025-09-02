package model_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/canonical/rt-conf/src/kcmd"
	"github.com/canonical/rt-conf/src/model"
)

const (
	grubSample = `
GRUB_DEFAULT=0
GRUB_TIMEOUT_STYLE=hidden
GRUB_TIMEOUT=0
GRUB_CMDLINE_LINUX_DEFAULT="quiet splash"
GRUB_CMDLINE_LINUX=""
`
)

type TestCase struct {
	Name        string
	Yaml        string
	Validations []struct {
		param string
		value string
	}
}

func mainLogic(t *testing.T, c TestCase, i int) (string, error) {
	dir := t.TempDir()
	// Set up temporary config-file
	tempConfigPath := filepath.Join(dir, fmt.Sprintf("config-%d.yaml", i))
	if err := os.WriteFile(tempConfigPath, []byte(c.Yaml), 0o644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	// Setup temporary custom grub config file
	tempCustomCfgPath := filepath.Join(dir, fmt.Sprintf("rt-conf-%d.cfg", i))
	if err := os.WriteFile(tempCustomCfgPath, []byte(grubSample), 0o644); err != nil {
		t.Fatalf("failed to write grub: %v", err)
	}

	t.Cleanup(func() {
		if err := os.Remove(tempConfigPath); err != nil {
			t.Fatal(err)
		}
		if err := os.Remove(tempCustomCfgPath); err != nil {
			t.Fatal(err)
		}
	})

	t.Logf("tempConfigPath: %s\n", tempConfigPath)
	t.Logf("tempCustomCfgPath: %s\n", tempCustomCfgPath)

	var conf model.InternalConfig
	if err := conf.Data.LoadFromFile(tempConfigPath); err != nil {
		return "", fmt.Errorf("failed to load config file: %v", err)
	}

	conf.GrubCfg = model.Grub{
		CustomGrubFilePath: tempCustomCfgPath,
	}

	t.Logf("Config: %+v\n", conf)

	// Run the InjectToFile method
	_, err := kcmd.ProcessKcmdArgs(&conf)
	if err != nil {
		return "", fmt.Errorf("ProcessKcmdArgs failed: %v", err)
	}

	// Verify default-grub updates
	updatedGrub, err := os.ReadFile(tempCustomCfgPath)
	if err != nil {
		return "", fmt.Errorf("failed to read modified grub file: %v", err)
	}

	t.Log("\nGrub file: \n", string(updatedGrub))

	return string(updatedGrub), nil
}

func TestHappyYamlKcmd(t *testing.T) {
	happyCases := []TestCase{
		{
			Name: "Using all cpus 0-N",
			Yaml: `
kernel-cmdline:
  parameters:
    - isolcpus=0-N
    - nohz=on
    - nohz_full=0-N
    - kthread_cpus=0-N
    - irqaffinity=0-N
`,
			Validations: []struct {
				param string
				value string
			}{
				{"isolcpus", "0-N"},
				{"nohz", "on"},
				{"nohz_full", "0-N"},
				{"kthread_cpus", "0-N"},
				{"irqaffinity", "0-N"},
			},
		},
		{
			Name: "Isolating cpu 0",
			Yaml: `
kernel-cmdline:
  parameters:
    - isolcpus=0
    - nohz=off
    - nohz_full=1-N
    - kthread_cpus=1-N
    - irqaffinity=1-N
`,
			Validations: []struct {
				param string
				value string
			}{
				{"isolcpus", "0"},
				{"nohz", "off"},
				{"nohz_full", "1-N"},
				{"kthread_cpus", "1-N"},
				{"irqaffinity", "1-N"},
			},
		},
	}

	for i, c := range happyCases {
		s, err := mainLogic(t, c, i)
		if err != nil {
			t.Fatal(err)
		}
		t.Run("HappyCases", func(t *testing.T) {
			for j, tc := range c.Validations {
				t.Log("Test case: ", j)
				if !strings.Contains(s,
					fmt.Sprintf("%s=%s", tc.param, tc.value)) {
					t.Errorf("\nExpected %s=%s in grub file, but not found",
						tc.param, tc.value)
				}
			}
		})
	}
}

func TestUnhappyYamlKcmd(t *testing.T) {
	// Building long kernel parameters list for testing
	// case that violates COMMAND_LINE_SIZE limit
	var sb strings.Builder
	sb.WriteString("kernel-cmdline:\n")
	// If each key=value parameter is 31 characters long
	// So being x the maximum number of parameters:
	// 31 * x + x <= 2048 (+x because we need to count the spaces between parameters)
	// The maximum x is 64 so with 65 parameters we exceed the limit
	for i := range 65 {
		sb.WriteString(fmt.Sprintf("  - some_kernel_cmd_option%03d=value\n", i))
	}
	longYamlList := sb.String()

	UnhappyCases := []TestCase{
		{
			Name:        "Command line exceeds maximum length",
			Yaml:        longYamlList,
			Validations: nil,
		},
		{
			Name: "Invalid parameter name",
			Yaml: `
kernel-cmdline:
  parameters:
		- 34foo=a
`,
			Validations: nil,
		},
		{
			Name: "isolcpus a is invalid",
			Yaml: `
kernel-cmdline:
  parameters:
		- isolcpus=a
		- nohz=on
		- nohz_full=0-N
		- kthread_cpus=0
		- irqaffinity=0-N
`,
			Validations: nil,
		},
		{
			Name: "kthread_cpus z is invalid",
			Yaml: `
kernel-cmdline:
  parameters:
		- isolcpus=0
		- nohz=on
		- nohz_full=0-N
		- kthread_cpus=z
		- irqaffinity=0-N
`,
			Validations: nil,
		},
		{
			Name: "nohz true is invalid",
			Yaml: `
kernel-cmdline:
  parameters:
		- isolcpus=0
		- nohz=true
		- nohz_full=0-N
		- kthread_cpus=0-N
		- irqaffinity=0-N
`,
			Validations: nil,
		},
		{
			Name: "isolcpus 100000000 is invalid",
			Yaml: `
kernel-cmdline:
  parameters:
		- isolcpus=100000000
		- nohz=off
		- nohz_full=0-N
		- kthread_cpus=0-N
		- irqaffinity=0-N
`,
			Validations: nil,
		},
	}

	for i, c := range UnhappyCases {
		t.Run(c.Name, func(t *testing.T) {
			_, err := mainLogic(t, c, i)
			if err == nil {
				t.Fatalf("Expected error, but got nil on YAML: \n%v",
					c.Yaml)
			}
		})
	}
}
