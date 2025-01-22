package cpu

import (
	"testing"
)

func TestComplement(t *testing.T) {
	type test struct {
		input  string
		output string
		tCores int
	}

	var tst = []test{
		{
			"0-7",
			"8,9",
			10,
		},
		{
			"0-1",
			"2,3,4,5,6,7,8,9",
			10,
		},
		{
			"2",
			"0,1,3",
			4,
		},
		{
			"0-20:2/5", // 0,1,5,6,10,11,15,16,20
			"2,3,4,7,8,9,12,13,14,17,18,19,21,22,23",
			24,
		},
		{
			"2,0-20:2/5,21-23", // 0,1,2,5,6,10,11,15,16,20,21,22,23
			"3,4,7,8,9,12,13,14,17,18,19",
			24,
		},
	}

	for _, tt := range tst {
		t.Run(tt.input, func(t *testing.T) {
			res, err := GenerateComplementCPUList(tt.input, tt.tCores)
			if err != nil {
				t.Fatalf("Failed GenerateComplementCPUList: %v", err)
			}
			if res != tt.output {
				t.Fatalf("expected %v, got %v", tt.output, res)
			}
		})
	}
}

func TestMutualExclusionCheck(t *testing.T) {
	type test struct {
		s1    string
		s2    string
		ncpus int
		mutE  bool // expected result
	}

	var tst = []test{
		{
			"0-7",
			"8-9",
			10,
			true,
		},
		{
			"1-99:2/10",
			"3-99:3/10",
			100,
			true,
		},
		{
			"1,2,3",
			"5,6",
			8,
			true,
		},
		{
			"1,2,3",
			"3,4",
			8,
			false,
		},
		{
			"0-4",
			"3,7",
			8,
			false,
		},
	}
	for _, tt := range tst {
		t.Run(tt.s1, func(t *testing.T) {

			check, err := cpuListsExclusive(tt.s1,
				tt.s2, tt.ncpus)
			if err != nil {
				t.Fatalf("MutuallyExclusive check failed: %v", err)
			}
			if check != tt.mutE {
				t.Logf("Test: %v", tt)
				t.Fatalf("expected %v, got %v", tt.mutE, check)
			}
		})
	}

}
