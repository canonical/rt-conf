package kcmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
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
