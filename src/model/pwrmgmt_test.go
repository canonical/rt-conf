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
