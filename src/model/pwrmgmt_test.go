package model

import "testing"

func TestPwrMgmtValidationHappy(t *testing.T) {
	var happyCases = []CpuGovernanceRule{
		{
			CPUs:    "0",
			ScalGov: "balanced",
		},
		{
			CPUs:    "0",
			ScalGov: "powersave",
		},
		{
			CPUs:    "0",
			ScalGov: "performance",
		},
	}
	for i, tc := range happyCases {
		t.Run("case-"+string(rune(i)), func(t *testing.T) {
			err := tc.Validate()
			if err != nil {
				t.Fatalf("error: %v", err)
			}
		})
	}
}

func TestPwrMgmtValidationUnhappy(t *testing.T) {
	var happyCases = []struct {
		err    string
		sclgov CpuGovernanceRule
	}{
		{
			"invalid cpu scaling governor: perf",
			CpuGovernanceRule{
				CPUs:    "0",
				ScalGov: "perf",
			},
		},
		{
			"invalid cpu scaling governor: pwer",
			CpuGovernanceRule{
				CPUs:    "0",
				ScalGov: "pwer",
			},
		},
		{
			"invalid cpu scaling governor: balance",
			CpuGovernanceRule{
				CPUs:    "0",
				ScalGov: "balance",
			},
		},
	}
	for i, tc := range happyCases {
		t.Run("case-"+string(rune(i)), func(t *testing.T) {
			err := tc.sclgov.Validate()
			if err == nil {
				t.Fatalf("Expected error on test #%v: %v", i, tc.sclgov)
			}
			if err.Error() != tc.err {
				t.Fatalf("Expected error message: %s, got: %s", tc.err, err)
			}

		})
	}
}

func TestCheckFreqFormat(t *testing.T) {
	tests := []struct {
		name    string
		rule    CpuGovernanceRule
		wantErr string
	}{
		{
			name: "valid GHz and MHz",
			rule: CpuGovernanceRule{
				MinFreq: "1.8GHz",
				MaxFreq: "3.4GHz",
			},
			wantErr: "",
		},
		{
			name: "valid G and M suffix only",
			rule: CpuGovernanceRule{
				MinFreq: "1.2M",
				MaxFreq: "2.4G",
			},
			wantErr: "",
		},
		{
			name: "valid no unit suffix",
			rule: CpuGovernanceRule{
				MinFreq: "1800000",
				MaxFreq: "2400000",
			},
			wantErr: "",
		},
		{
			name: "valid float without unit",
			rule: CpuGovernanceRule{
				MinFreq: "1.2",
				MaxFreq: "2.5",
			},
			wantErr: "",
		},
		{
			name: "valid lowercase hz",
			rule: CpuGovernanceRule{
				MinFreq: "1.0mhz",
				MaxFreq: "3.0gghz",
			},
			wantErr: "invalid max frequency: invalid frequency format: 3.0gghz",
		},
		{
			name: "invalid max freq string",
			rule: CpuGovernanceRule{
				MinFreq: "1.2GHz",
				MaxFreq: "threeGHz",
			},
			wantErr: "invalid max frequency: invalid frequency format: threeGHz",
		},
		{
			name: "invalid min freq string",
			rule: CpuGovernanceRule{
				MinFreq: "oneGHz",
				MaxFreq: "3.4GHz",
			},
			wantErr: "invalid min frequency: invalid frequency format: oneGHz",
		},
		{
			name: "both max and min invalid",
			rule: CpuGovernanceRule{
				MinFreq: "badMin",
				MaxFreq: "badMax",
			},
			wantErr: "invalid min frequency: invalid frequency format: badMin",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.rule.CheckFreqFormat()
			if tc.wantErr == "" && err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error: %q, got nil", tc.wantErr)
				}
				if err.Error() != tc.wantErr {
					t.Errorf("expected error: %q, got: %q",
						tc.wantErr,
						err.Error())
				}
			}
		})
	}
}
