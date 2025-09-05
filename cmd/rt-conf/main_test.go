package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunHappy(t *testing.T) {
	type tests struct {
		name string
		args []string
		yaml string
	}

	tmpdir := t.TempDir()
	configPath := filepath.Join(tmpdir, "config.yaml")

	testCases := []tests{
		{
			name: "Valid empty config",
			args: []string{"rt-conf", "-file", configPath},
			yaml: `
kernel-cmdline:
cpu-governance:
irq-tuning:
`,
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			if test.yaml != "" {
				if err := os.WriteFile(configPath,
					[]byte(test.yaml), 0o644); err != nil {
					t.Fatalf("failed to write file: %v", err)
				}
			}
			err := run(test.args)
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
		})
	}
}

func TestRunUnhappy(t *testing.T) {
	type envVar struct {
		key   string
		value string
	}

	type tests struct {
		name string // test name
		args []string
		err  string
		yaml string
		envs []envVar
	}

	tmpdir := t.TempDir()
	configPath := filepath.Join(tmpdir, "config.yaml")

	testCases := []tests{
		{
			name: "Invalid VERBOSE value",
			args: []string{"rt-conf", "-file", configPath},
			envs: []envVar{{key: "VERBOSE", value: "yes"}},
			err:  "failed to parse VERBOSE env var",
			yaml: `
kernel-cmdline:
cpu-governance:
irq-tuning:
`,
		},
		{
			name: "No config path",
			args: []string{"rt-conf"},
			err:  "path not set",
			yaml: "",
		},
		{
			name: "Empty config | All commented out",
			args: []string{"rt-conf", "-file", configPath},
			err:  "failed to load config file: empty config file",
			yaml: `
# Kernel command line parameters
`,
		},
		{
			name: "Invalid config path",
			args: []string{"rt-conf", "-file", "/does/not/exist"},
			err:  "failed to find file",
			yaml: "",
		},
		{
			name: "Invalid config + verbose mode",
			args: []string{"rt-conf", "--verbose", "-file", "/does/not/exist"},
			err:  "failed to find file",
			yaml: "",
		},
		{
			name: "Error processing kernel cmdline args",
			args: []string{"rt-conf", "-file", configPath},
			err:  "failed to process kernel cmdline args",
			yaml: `
kernel-cmdline:
  parameters:
    - nohz=on
`,
		},
		{
			name: "No IRQs found",
			args: []string{"rt-conf", "-file", configPath},
			err:  "failed to process interrupts",
			yaml: `
irq-tuning:
  "foo":
    cpus: "0"
    filter:
      actions: "xxxxxx"
`,
		},
		{
			name: "Failed to process power management config",
			args: []string{"rt-conf", "-file", configPath},
			err:  "failed to process power management config",
			yaml: `
cpu-governance:
  "bar":
    cpus: "0"
    scaling-governor: "performance"
`,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			if len(test.envs) > 0 {
				for _, env := range test.envs {
					if err := os.Setenv(env.key, env.value); err != nil {
						t.Fatalf("failed to set env: %v", err)
					}
				}
			}

			if test.yaml != "" {
				if err := os.WriteFile(configPath,
					[]byte(test.yaml), 0o644); err != nil {
					t.Fatalf("failed to write file: %v", err)
				}
			}
			err := run(test.args)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !strings.Contains(err.Error(), test.err) {
				t.Fatalf("expected error '%s', got: '%v'", test.err, err)
			}

			t.Cleanup(func() {
				if len(test.envs) > 0 {
					for _, env := range test.envs {
						os.Unsetenv(env.key)
					}
				}
			})
		})
	}
}
