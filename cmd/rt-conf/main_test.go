package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/canonical/rt-conf/src/data"
	"github.com/canonical/rt-conf/src/interrupts"
	"github.com/canonical/rt-conf/src/kcmd"
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
	Yaml        string
	Validations []struct {
		param string
		value string
	}
}

type IRQTestCase struct {
	Yaml   string
	Writer interrupts.IRQWriter
	Reader interrupts.IRQReader
}

func TestHappyIRQtunning(t *testing.T) {

	var happyCases = []IRQTestCase{
		{
			Yaml: `
irq_tunning:
- cpus: 0
  filter:
    number: 10
`,
			Writer: &interrupts.MockIRQWriter{},
			Reader: &interrupts.MockIRQReader{
				IRQs: map[uint]interrupts.IRQInfo{
					10: {
						Number: 10,
					},
				},
			},
		},
	}

	for i, c := range happyCases {
		t.Run("Happy Cases", func(t *testing.T) {
			_, err := mainLogicIRQ(t, c, i)
			if err != nil {
				t.Fatalf("On YAML: \n%v\nError: %v", c.Yaml, err)
			}
		})
	}
}

func TestUnhappyIRQtunning(t *testing.T) {

	var UnhappyCases = []IRQTestCase{
		{
			// Invalid number
			Yaml: `
irq_tunning:
- cpus: 0
  filter:
    number: a
`,
			Writer: &interrupts.MockIRQWriter{},
			Reader: &interrupts.MockIRQReader{},
		},
		{
			// Invalid RegEx
			Yaml: `
irq_tunning:
- cpus: 0
  filter:
    number: 0
    action: "*"
`,
			Writer: &interrupts.MockIRQWriter{},
			Reader: &interrupts.MockIRQReader{},
		},
	}

	for i, c := range UnhappyCases {
		t.Run("Unhappy Cases", func(t *testing.T) {
			_, err := mainLogicIRQ(t, c, i)
			// if err != nil {
			// 	t.Fatalf("On YAML: \n%v\nError: %v", c.Yaml, err)
			// }
			if err == nil {
				t.Fatalf("Expected error, got nil on YAML %v", c.Yaml)
			}
		})
	}
}

func mainLogicIRQ(t *testing.T, c IRQTestCase, i int) (string, error) {
	fmt.Println("[DEBUG] mainLogicIRQ")
	tempConfigPath := setupTempFile(t, c.Yaml, i)
	t.Cleanup(func() {
		os.Remove(tempConfigPath)
	})
	var conf data.InternalConfig
	if d, err := data.LoadConfigFile(tempConfigPath); err != nil {
		return "", fmt.Errorf("failed to load config file: %v", err)
	} else {
		conf.Data = *d
	}

	err := interrupts.ApplyIRQConfig(&conf, c.Reader, c.Writer)
	if err != nil {
		return "", fmt.Errorf("Failed to process interrupts: %v", err)
	}
	return "", nil
}

func setupTempFile(t *testing.T, content string, idex int) string {
	t.Helper()

	tmpFile, err := os.CreateTemp("", fmt.Sprintf("tempfile-%d", idex))
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	return tmpFile.Name()
}

func mainLogic(t *testing.T, c TestCase, i int) (string, error) {

	// Set up temporary config-file, grub.cfg and default-grub files
	tempConfigPath := setupTempFile(t, c.Yaml, i)
	tempGrubPath := setupTempFile(t, grubSample, i)
	t.Cleanup(func() {
		os.Remove(tempConfigPath)
		os.Remove(tempGrubPath)
	})

	t.Logf("tempConfigPath: %s\n", tempConfigPath)
	t.Logf("tempGrubPath: %s\n", tempGrubPath)

	var conf data.InternalConfig
	if d, err := data.LoadConfigFile(tempConfigPath); err != nil {
		return "", fmt.Errorf("failed to load config file: %v", err)
	} else {
		conf.Data = *d
	}

	conf.CfgFile = tempConfigPath
	conf.GrubDefault = data.Grub{
		File:    tempGrubPath,
		Pattern: data.PatternGrubDefault,
	}

	t.Logf("Config: %+v\n", conf)

	// Run the InjectToFile method
	_, err := kcmd.ProcessKcmdArgs(&conf)
	if err != nil {
		return "", fmt.Errorf("ProcessKcmdArgs failed: %v", err)
	}

	// Verify default-grub updates
	updatedGrub, err := os.ReadFile(tempGrubPath)
	if err != nil {
		return "", fmt.Errorf("failed to read modified grub file: %v", err)
	}

	t.Log("\nGrub file: ", string(updatedGrub))

	return string(updatedGrub), nil
}

func TestHappyMainLogic(t *testing.T) {
	var happyCases = []TestCase{
		{
			Yaml: `
kernel_cmdline:
  isolcpus: "0-n"
  nohz: "on"
  nohz_full: "0-n"
  kthread_cpus: "0-n"
  irqaffinity: "0-n"
`,
			Validations: []struct {
				param string
				value string
			}{
				{"isolcpus", "0-n"},
				{"nohz", "on"},
				{"nohz_full", "0-n"},
				{"kthread_cpus", "0-n"},
				{"irqaffinity", "0-n"},
			},
		},
		{
			Yaml: `
kernel_cmdline:
  isolcpus: "0"
  nohz: "off"
  nohz_full: "0-n"
  kthread_cpus: "0-n"
  irqaffinity: "0-n"
`,
			Validations: []struct {
				param string
				value string
			}{
				{"isolcpus", "0"},
				{"nohz", "off"},
				{"nohz_full", "0-n"},
				{"kthread_cpus", "0-n"},
				{"irqaffinity", "0-n"},
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

func TestUnhappyMainLogic(t *testing.T) {
	var UnhappyCases = []TestCase{
		{
			// isolcpus: "a" is valid
			Yaml: `
kernel_cmdline:
  isolcpus: "a"
  nohz: "on"
  nohz_full: "0-n"
  kthread_cpus: "0"
  irqaffinity: "0-n"
`,
			Validations: nil,
		},
		{
			// irqaffinity: "z" is valid
			Yaml: `
kernel_cmdline:
  isolcpus: "0"
  nohz: "on"
  nohz_full: "0-n"
  kthread_cpus: "z"
  irqaffinity: "0-n"
`,
			Validations: nil,
		},
		{
			// nohz: "true" is valid it should be 'on' or 'off'
			Yaml: `
kernel_cmdline:
  isolcpus: "0"
  nohz: "true"
  nohz_full: "0-n"
  kthread_cpus: "0-n"
  irqaffinity: "0-n"
`,
			Validations: nil,
		},
		{
			// isolcpus: "100000000" is invalid
			Yaml: `
kernel_cmdline:
  isolcpus: "100000000"
  nohz: "off"
  nohz_full: "0-n"
  kthread_cpus: "0-n"
  irqaffinity: "0-n"
`,
			Validations: nil,
		},
	}
	for i, c := range UnhappyCases {
		t.Run("UnhappyCases", func(t *testing.T) {
			_, err := mainLogic(t, c, i)
			if err == nil {
				t.Fatalf("Expected error, but got nil on YAML: \n%v",
					c.Yaml)
			}
		})
	}
}
