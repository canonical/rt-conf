package kcmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/canonical/rt-conf/src/model"
)

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

func TestParseGrubFileHappy(t *testing.T) {
	var testCases = []struct {
		content string
		params  map[string]string
	}{
		{
			content: `GRUB_HIDDEN_TIMEOUT_QUIET=true`,
			params: map[string]string{
				"GRUB_HIDDEN_TIMEOUT_QUIET": "true",
			},
		},
		{
			content: `GRUB_TIMEOUT=2`,
			params: map[string]string{
				"GRUB_TIMEOUT": "2",
			},
		},
		{
			content: `GRUB_CMDLINE_LINUX_DEFAULT="rootfstype=ext4 quiet splash acpi_osi="`,
			params: map[string]string{
				"GRUB_CMDLINE_LINUX_DEFAULT": "rootfstype=ext4 quiet splash acpi_osi=",
			},
		},
		{
			content: "GRUB_DEFAULT=0\n" +
				"#GRUB_HIDDEN_TIMEOUT=0\n" +
				"GRUB_HIDDEN_TIMEOUT_QUIET=true\n" +
				"GRUB_TIMEOUT=2\n" +
				"GRUB_DISTRIBUTOR=`lsb_release -i -s 2> /dev/null || echo Debian`\n" +
				"GRUB_CMDLINE_LINUX_DEFAULT=\"rootfstype=ext4 quiet splash acpi_osi=\"\n" +
				"GRUB_CMDLINE_LINUX=\"\"\n",

			params: map[string]string{
				"GRUB_DEFAULT":               "0",
				"GRUB_HIDDEN_TIMEOUT_QUIET":  "true",
				"GRUB_TIMEOUT":               "2",
				"GRUB_DISTRIBUTOR":           "`lsb_release -i -s 2> /dev/null || echo Debian`",
				"GRUB_CMDLINE_LINUX_DEFAULT": "rootfstype=ext4 quiet splash acpi_osi=",
				"GRUB_CMDLINE_LINUX":         "",
			},
		},
	}
	for i, tc := range testCases {
		tmpFile := setupTempFile(t, tc.content, i)
		t.Cleanup(func() {
			os.Remove(tmpFile)
		})

		params, err := ParseDefaultGrubFile(tmpFile)
		if err != nil {
			t.Fatalf("Failed to parse grub file: %v", err)
		}
		for k, v := range params {
			vt, ok := tc.params[k]
			if !ok {
				t.Fatalf("Expected %s not found", k)
			}
			if v != vt {
				t.Fatalf("Expected %s=%s, got %s=%s", k, vt, k, v)
			}

		}
	}
}

func TestDuplicatedParams(t *testing.T) {
	var testCases = []struct {
		name    string
		cmdline string
		err     error
	}{
		{
			name:    "No duplicates",
			cmdline: "quiet splash foo",
			err:     nil,
		},
		{
			name:    "Single parameter",
			cmdline: "quiet",
			err:     nil,
		},
		{
			name:    "Duplicate boolean parameters",
			cmdline: "quiet splash quiet",
			err:     nil,
		},
		{
			name:    "Duplicate keys with different values",
			cmdline: "potato=mashed potato=salad",
			err:     errors.New("duplicated parameter:"),
		},
		{
			name:    "Duplicate key-value pairs",
			cmdline: "potato=pie potato=pie",
			err:     nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := duplicatedParams(tc.cmdline)
			if tc.err != nil {
				if err == nil {
					t.Fatalf("Expected error %v, got nil", tc.err)
				}
				if !strings.Contains(err.Error(), tc.err.Error()) {
					t.Fatalf("Expected error %v, got %v", tc.err, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
		})
	}
}

func TestUpdateGrub(t *testing.T) {
	tests := []struct {
		name         string
		grubContent  string
		kcmd         model.KernelCmdline
		expectErr    string
		expectOutput string
	}{
		{
			name:        "No params to inject",
			grubContent: ``,
			kcmd:        model.KernelCmdline{},
			expectErr:   "no parameters to inject",
		},
		{
			name:        "ParseDefaultGrubFile fails",
			grubContent: "", // file will be removed
			kcmd: model.KernelCmdline{
				IsolCPUs: "1-3",
			},
			expectErr: "failed to parse grub file",
		},
		{
			name:        "GRUB_CMDLINE_LINUX missing",
			grubContent: `GRUB_TIMEOUT=5`,
			kcmd: model.KernelCmdline{
				IsolCPUs: "1-3",
			},
			expectErr: "GRUB_CMDLINE_LINUX not found",
		},
		{
			name: "Duplicate params found",
			grubContent: `GRUB_CMDLINE_LINUX="isolcpus=1-3 isolcpus=2-4"
`,
			kcmd: model.KernelCmdline{
				IsolCPUs: "2-4",
			},
			expectErr: "invalid existing parameters",
		},
		{
			name: "ProcessFile fails",
			grubContent: `GRUB_CMDLINE_LINUX="isolcpus=1-3"
`,
			kcmd: model.KernelCmdline{
				Nohz: "on",
			},
			expectErr: "error updating",
		},
		{
			name: "Success",
			grubContent: `GRUB_CMDLINE_LINUX="isolcpus=1-3"
`,
			kcmd: model.KernelCmdline{
				IsolCPUs: "1-3",
				Nohz:     "on",
			},
			expectOutput: "Detected bootloader: GRUB",
		},
	}

	model.PatternGrubDefault = regexp.MustCompile(`^(GRUB_CMDLINE_LINUX=")([^"]*)(")$`)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			grubPath := filepath.Join(tmpDir, "grub")

			// If grub content exists, write it
			if tc.grubContent != "" {
				if err := os.WriteFile(grubPath, []byte(tc.grubContent), 0644); err != nil {
					t.Fatal(err)
				}
			}

			// If test needs Parse failure, remove the file
			if tc.name == "ParseDefaultGrubFile fails" {
				os.Remove(grubPath)
			}

			conf := &model.InternalConfig{
				Data: model.Config{
					KernelCmdline: tc.kcmd,
				},
				GrubDefault: model.Grub{
					File: grubPath,
				},
			}

			// Patch ProcessFile to simulate failure
			origProcessFile := model.ProcessFile
			defer func() { model.ProcessFile = origProcessFile }()
			if strings.Contains(tc.name, "ProcessFile fails") {
				model.ProcessFile = func(tf model.FileTransformer) error {
					return fmt.Errorf("mock write failure")
				}
			} else {
				model.ProcessFile = func(tf model.FileTransformer) error {
					// simulate a successful file process
					return nil
				}
			}

			msgs, err := UpdateGrub(conf)
			if tc.expectErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.expectErr) {
					t.Fatalf("expected error %q, got: %v", tc.expectErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				found := false
				for _, msg := range msgs {
					if strings.Contains(msg, tc.expectOutput) {
						found = true
						break
					}
				}
				if !found {
					t.Fatalf("expected output to contain %q, got %v", tc.expectOutput, msgs)
				}
			}
		})
	}
}
