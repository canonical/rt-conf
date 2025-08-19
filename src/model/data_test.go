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
					"isolcpus=1",
					"nohz=on",
					"nohz_full=1",
					"kthread_cpus=0",
					"irqaffinity=0",
				},
			},
			err: nil,
		},
		{
			name: "Invalid kernel cmdline",
			cfg: &Config{
				KernelCmdline: KernelCmdline{
					"nohz=potato", // Invalid value
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
