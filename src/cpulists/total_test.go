package cpulists

import (
	"errors"
	"os/exec"
	"strings"
	"testing"
)

func TestTotalCPUs(t *testing.T) {
	c, err := totalCPUs()
	if err != nil {
		t.Fatalf("Failed to get total CPUs: %s", err)
	}
	if c == 0 {
		t.Fatal("Unexpected total CPUs: 0")
	}
}

func TestTotalCPUsUnhappy(t *testing.T) {
	type testCase struct {
		name       string
		mockOutput []byte
		mockErr    error
		expectErr  string
	}

	var testCases = []testCase{
		{
			name:       "command error",
			mockOutput: nil,
			mockErr:    errors.New("fake command error"),
			expectErr:  "fake command error",
		},
		{
			name:       "invalid json output",
			mockOutput: []byte("not-json"),
			expectErr:  "invalid character",
		},
		{
			name: "invalid CPU number",
			mockOutput: []byte(`{
				"lscpu": [
					{ "field": "CPU(s):", "data": "not-an-int", "children": [] }
				]
			}`),
			expectErr: "invalid syntax",
		},
		{
			name: "missing CPU field",
			mockOutput: []byte(`{
				"lscpu": [
					{ "field": "Model name:", "data": "Intel Xeon", "children": [] }
				]
			}`),
			expectErr: "could not find total CPUs",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// override the execCommand used by totalCPUs
			execCommand = func(name string, args ...string) ([]byte, error) {
				return tc.mockOutput, tc.mockErr
			}
			defer func() {
				execCommand = func(name string, args ...string) ([]byte, error) {
					return exec.Command(name, args...).Output()
				}
			}()

			_, err := totalCPUs()
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.expectErr) {
				t.Fatalf("expected error containing %q, got: %v", tc.expectErr, err)
			}
		})
	}
}
