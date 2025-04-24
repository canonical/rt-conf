package model

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestIsOwnedByRoot(t *testing.T) {
	tmpdir := t.TempDir()
	filePath := filepath.Join(tmpdir, "testfile")
	err := os.WriteFile(filePath,
		[]byte("test"), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	fi, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("failed to stat test file: %v", err)
	}
	if isRoot := IsOwnedByRoot(fi); isRoot {
		t.Fatalf("expected false, got true")
	}

}

func TestLoadConfigFile(t *testing.T) {
	var testCases = []struct {
		name        string
		yaml        string
		perm        os.FileMode
		cfg         *Config
		err         error
		ownedByRoot bool
	}{
		{
			name:        "FileDoesNotExist",
			yaml:        "",
			perm:        0644,
			cfg:         nil,
			err:         fmt.Errorf("failed to find file"),
			ownedByRoot: true,
		},
		{
			name:        "FileInvalidPermissions",
			yaml:        "kernel_cmdline:",
			perm:        0755,
			cfg:         nil,
			err:         fmt.Errorf("has invalid permissions"),
			ownedByRoot: true,
		},
		{
			name:        "FileNotOwnedByRoot",
			yaml:        `kernel_cmdline:`,
			perm:        0644,
			cfg:         nil,
			err:         fmt.Errorf("not owned by root"),
			ownedByRoot: false,
		},
		{
			name:        "FailedToUnmarshalYAML",
			yaml:        `kernel_cmdline: {`,
			perm:        0644,
			cfg:         nil,
			err:         fmt.Errorf("failed to unmarshal data"),
			ownedByRoot: true,
		},
		{
			name: "FileValid",
			yaml: `
kernel_cmdline:
    nohz: "on"
    isolcpus: "1"
    kthread_cpus: "0"
    irqaffinity: "0"
`,
			perm: 0644,
			cfg: &Config{
				KernelCmdline: KernelCmdline{
					Nohz:        "on",
					IsolCPUs:    "1",
					KthreadCPUs: "0",
					IRQaffinity: "0",
				},
			},
			ownedByRoot: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpdir := t.TempDir()
			cfgFilePath := filepath.Join(tmpdir, "config.yaml")

			IsOwnedByRoot = func(_ os.FileInfo) bool {
				return tc.ownedByRoot
			}

			if tc.yaml != "" {
				if err := os.WriteFile(cfgFilePath,
					[]byte(tc.yaml), tc.perm); err != nil {
					t.Fatalf("failed to write file: %v", err)
				}
			}

			cfg, err := LoadConfigFile(cfgFilePath)

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
