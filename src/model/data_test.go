package model

import (
	"errors"
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name string
		cfg  *Config
		err  error
	}{
		{
			name: "Empty config",
			cfg:  &Config{},
			err:  nil,
		},
		{
			name: "Valid config",
			cfg: &Config{
				KernelCmdline: KernelCmdline{
					IsolCPUs:    "1",
					Nohz:        "on",
					NohzFull:    "1",
					KthreadCPUs: "0",
					IRQaffinity: "0",
				},
			},
			err: nil,
		},
		{
			name: "Invalid kernel cmdline",
			cfg: &Config{
				KernelCmdline: KernelCmdline{
					Nohz: "potato", // Invalid value
				},
			},
			err: errors.New("failed to validate kernel cmdline"),
		},
		{
			name: "Invalid IRQ tuning",
			cfg: &Config{
				Interrupts: Interrupts{
					"invalid": {
						CPUs: "0",
						Filter: IRQFilter{
							Name: "**", // Invalid regex
						},
					},
				},
			},
			err: errors.New("failed to validate irq tuning"),
		},
		{
			name: "Invalid name for IRQ tuning rule - with space",
			cfg: &Config{
				Interrupts: Interrupts{
					"foo bar buzz": {
						CPUs: "0",
						Filter: IRQFilter{
							Name: "foo",
						},
					},
				},
			},
			err: errors.New("rule name cannot contain whitespace characters"),
		},
		{
			name: "Invalid CPU governance rule",
			cfg: &Config{
				CpuGovernance: PwrMgmt{
					"foo": {
						CPUs:    "0-1",
						ScalGov: "potato", // Invalid governor
					},
				},
			},
			err: errors.New("failed to validate cpu governance rule"),
		},
		{
			name: "Invalid name for CPU governance rule - with space",
			cfg: &Config{
				CpuGovernance: PwrMgmt{
					"foo bar buzz": {
						CPUs:    "0-1",
						ScalGov: "powersave",
					},
				},
			},
			err: errors.New("rule name cannot contain whitespace characters"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.cfg.Validate()
			// Unhappy cases
			if tc.err != nil {
				if err == nil {
					t.Fatalf("expected error '%v', got nil", tc.err)
				}
				if !strings.Contains(err.Error(), tc.err.Error()) {
					t.Fatalf("expected error '%v', got '%v'", tc.err, err)
				}
				return
			}
			// Happy cases
			if err != nil {
				t.Fatalf("expected no error, got '%v'", err)
			}
		})
	}
}
