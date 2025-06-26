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
				MaxFreq: "3.4GHz",
				MinFreq: "1.8GHz",
			},
			wantErr: "",
		},
		{
			name: "valid G and M suffix only",
			rule: CpuGovernanceRule{
				MaxFreq: "2.4G",
				MinFreq: "1.2M",
			},
			wantErr: "",
		},
		{
			name: "valid no unit suffix",
			rule: CpuGovernanceRule{
				MaxFreq: "2400000",
				MinFreq: "1800000",
			},
			wantErr: "",
		},
		{
			name: "valid float without unit",
			rule: CpuGovernanceRule{
				MaxFreq: "2.5",
				MinFreq: "1.2",
			},
			wantErr: "",
		},
		{
			name: "valid lowercase hz",
			rule: CpuGovernanceRule{
				MaxFreq: "3.0ghz",
				MinFreq: "1.0mhz",
			},
			wantErr: "invalid frequency format: 3.0ghz",
		},
		{
			name: "invalid max freq string",
			rule: CpuGovernanceRule{
				MaxFreq: "threeGHz",
				MinFreq: "1.2GHz",
			},
			wantErr: "invalid frequency format: threeGHz",
		},
		{
			name: "invalid min freq string",
			rule: CpuGovernanceRule{
				MaxFreq: "3.4GHz",
				MinFreq: "oneGHz",
			},
			wantErr: "invalid frequency format: oneGHz",
		},
		{
			name: "both max and min invalid",
			rule: CpuGovernanceRule{
				MaxFreq: "badMax",
				MinFreq: "badMin",
			},
			wantErr: "invalid frequency format: badMax", // Only first invalid is returned
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
					t.Errorf("expected error: %q, got: %q", tc.wantErr, err.Error())
				}
			}
		})
	}
}
