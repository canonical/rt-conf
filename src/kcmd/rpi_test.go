package kcmd

import (
	"strings"
	"testing"

	"github.com/canonical/rt-conf/src/model"
)

func TestUpdateRPi(t *testing.T) {
	type testCase struct {
		name        string
		kcmdline    model.KernelCmdline
		expectErr   string
		expectParts []string
	}

	tests := []testCase{
		{
			name:      "No parameters",
			kcmdline:  model.KernelCmdline{}, // all fields empty
			expectErr: "no parameters to inject",
		},
		{
			name: "Valid parameters",
			kcmdline: model.KernelCmdline{
				Parameters: []string{
					"isolcpus=1-3",
					"nohz=on",
					"nohz_full=1-3",
					"kthread_cpus=0",
					"irqaffinity=0",
				},
			},
			expectParts: []string{
				"isolcpus=1-3",
				"nohz=on",
				"nohz_full=1-3",
				"kthread_cpus=0",
				"irqaffinity=0",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &model.InternalConfig{
				Data: model.Config{
					KernelCmdline: tc.kcmdline,
				},
			}

			result, err := UpdateRPi(cfg)
			if tc.expectErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.expectErr) {
					t.Fatalf("expected error %q, got %v", tc.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(result) == 0 {
				t.Fatal("expected result, got none")
			}

			joined := strings.Join(result, " ")
			for _, part := range tc.expectParts {
				if !strings.Contains(joined, part) {
					t.Errorf("expected result to contain %q, got %q", part, joined)
				}
			}
		})
	}
}
