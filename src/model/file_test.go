package model

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestReadYAML(t *testing.T) {
	type testCase struct {
		name string
		yaml string
		cfg  *Config
		err  error
	}

	tmpdir := t.TempDir()
	testCases := []testCase{
		{
			name: "FileDoesNotExist",
			yaml: "",
			cfg:  nil,
			err:  errors.New("failed to read file"),
		},
		{
			name: "InvalidYAML",
			yaml: `: this is not valid yaml`,
			cfg:  nil,
			err:  errors.New("failed to unmarshal data"),
		},
		{
			name: "EmptyConfig",
			yaml: ``,
			cfg:  nil,
			err:  errors.New("empty config file"),
		},
		{
			name: "ValidationFails",
			yaml: `
kernel-cmdline:
  nohz: "invalid_value"
`,
			cfg: nil,
			err: errors.New("failed to validate kernel cmdline"),
		},
		{
			name: "Success",
			yaml: `
kernel-cmdline:
  nohz: "on"
`,
			cfg: &Config{
				KernelCmdline: KernelCmdline{
					Nohz: "on",
				},
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var filePath string

			if tc.name == "FileDoesNotExist" {
				// Intentionally not writing the file
				filePath = filepath.Join(tmpdir, "nonexistent.yaml")
			} else {
				filePath = filepath.Join(tmpdir, tc.name+".yaml")
				if err := os.WriteFile(filePath, []byte(tc.yaml), 0644); err != nil {
					t.Fatalf("failed to write test YAML file: %v", err)
				}
			}

			cfg, err := ReadYAML(filePath)
			// Unhappy cases
			if tc.err != nil {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tc.err)
				}
				if !strings.Contains(err.Error(), tc.err.Error()) {
					t.Fatalf("expected error to contain %q, got: %v", tc.err, err)
				}
				return
			}
			// Happy cases
			if err != nil {
				t.Fatalf("expected no error, got '%v'", err)
			}
			if cfg == nil {
				t.Fatalf("expected non-nil config, got nil")
			}
			if !reflect.DeepEqual(cfg, tc.cfg) {
				t.Fatalf("expected config %v, got %v", tc.cfg, cfg)
			}
		})
	}
}
