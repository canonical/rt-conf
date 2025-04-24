package kcmd

import (
	"errors"
	"strings"
	"testing"

	"github.com/canonical/rt-conf/src/model"
	"github.com/canonical/rt-conf/src/system"
)

func TestProcessKcmdArgs(t *testing.T) {
	type testCase struct {
		name      string
		input     model.InternalConfig
		mockSys   system.SystemType
		mockResp  []string
		mockErr   error
		expect    []string
		expectErr string
	}

	fake := func(msgs []string, err error) func(*model.InternalConfig) ([]string, error) {
		return func(_ *model.InternalConfig) ([]string, error) {
			return msgs, err
		}
	}

	tests := []testCase{
		{
			name:      "Empty KernelCmdline",
			input:     model.InternalConfig{},
			expect:    nil,
			expectErr: "",
		},
		{
			name: "Unsupported Bootloader",
			input: model.InternalConfig{
				Data: model.Config{
					KernelCmdline: model.KernelCmdline{Nohz: "on"},
				},
			},
			mockSys:   system.Unknown,
			expectErr: "unsupported bootloader",
		},
		{
			name: "Processor returns error",
			input: model.InternalConfig{
				Data: model.Config{
					KernelCmdline: model.KernelCmdline{Nohz: "on"},
				},
			},
			mockSys:   system.Grub,
			mockResp:  nil,
			mockErr:   errors.New("simulated failure"),
			expectErr: "simulated failure",
		},
		// TODO: fix this test case
		// {
		// 	name:      "Failed to detect system",
		// 	mockErr:   errors.New("failed to detect system"),
		// 	expectErr: "failed to detect system",
		// },
		{
			name: "Processor returns success",
			input: model.InternalConfig{
				Data: model.Config{
					KernelCmdline: model.KernelCmdline{Nohz: "on"},
				},
			},
			mockSys:  system.Grub,
			mockResp: []string{"ok"},
			expect:   []string{"ok"},
		},
	}

	savedDetect := system.DetectSystem
	savedKcmdSys := make(map[system.SystemType]func(*model.InternalConfig) ([]string, error))
	for k, v := range kcmdSys {
		savedKcmdSys[k] = v
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Mock DetectSystem
			system.DetectSystem = nil
			system.DetectSystem = func() (system.SystemType, error) {
				return tc.mockSys, tc.mockErr
			}

			// Replace processor if test case provides it
			if tc.mockSys != 0 {
				kcmdSys[tc.mockSys] = fake(tc.mockResp, tc.mockErr)
			}

			out, err := ProcessKcmdArgs(&tc.input)

			if tc.expectErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.expectErr) {
					t.Fatalf("expected error containing %q, got: %v", tc.expectErr, err)
				}
				return
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(out) != len(tc.expect) {
				t.Fatalf("expected %v, got %v", tc.expect, out)
			}
			for i := range out {
				if out[i] != tc.expect[i] {
					t.Fatalf("expected %v, got %v", tc.expect, out)
				}
			}
		})
	}

	// Restore
	system.DetectSystem = savedDetect
	kcmdSys = savedKcmdSys
}
