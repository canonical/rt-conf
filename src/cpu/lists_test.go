package cpu_test

import (
	"reflect"
	"testing"

	i "github.com/canonical/rt-conf/src/cpu"
)

func TestParseCPUSlistsHappy(t *testing.T) {
	type test struct {
		input  string
		tCores int
		output i.CPUs
	}

	var tst = []test{
		// single CPU
		{
			"4",
			8,
			i.CPUs{4: true},
		},
		// 3 single CPUs
		{
			"4,5,9",
			10,
			i.CPUs{4: true, 5: true, 9: true},
		},
		// CPU range
		{
			"0-7",
			10,
			i.CPUs{0: true, 1: true, 2: true, 3: true, 4: true,
				5: true, 6: true, 7: true},
		},
		// two CPUs ranges
		{
			"0-2,4-7",
			10,
			i.CPUs{0: true, 1: true, 2: true, 4: true,
				5: true, 6: true, 7: true},
		},
		// CPU range + single CPU
		{
			"0-2,3",
			4,
			i.CPUs{0: true, 1: true, 2: true, 3: true},
		},
		// Formated CPU list
		{
			"0-20:2/5",
			24,
			i.CPUs{0: true, 1: true, 5: true, 6: true,
				10: true, 11: true, 15: true, 16: true, 20: true},
		},
		// Formated CPU list + a single CPU
		{
			"0-20:2/5,23",
			24,
			i.CPUs{0: true, 1: true, 5: true, 6: true, 10: true,
				11: true, 15: true, 16: true, 20: true, 23: true},
		},
	}

	for _, tt := range tst {
		t.Run(tt.input, func(t *testing.T) {
			res, err := i.ParseCPUs(tt.input, tt.tCores)
			if err != nil {
				t.Fatalf("ParseCPUs failed: %v", err)
			}
			if len(res) != len(tt.output) {
				t.Fatalf("expected %v, got %v", tt.output, res)
			}
			if !reflect.DeepEqual(res, tt.output) {
				t.Fatalf("expected %v, got %v", tt.output, res)
			}
		})
	}
}

func TestParseCPUSlistsUnhappy(t *testing.T) {
	type test struct {
		input  string
		tCores int
		err    string
	}

	var tst = []test{
		{
			"4",
			4,
			"CPU greater than total CPUs: 4",
		},
		{
			"a",
			4,
			"invalid CPU: a",
		},
		{
			"1-2-3",
			4,
			"invalid range: 1-2-3",
		},
		{
			"a-2",
			4,
			"invalid start of range: a",
		},
		{
			"1-a",
			4,
			"invalid end of range: a",
		},
		{
			"6-8",
			8,
			"end of range greater than total CPUs: 6-8",
		},
		{
			"5-2",
			8,
			"start of range greater than end: 5-2",
		},
		{
			"0--2:",
			4,
			"invalid range: 0--2",
		},
		{
			"0-:2",
			4,
			"invalid end of range: ",
		},
		{
			"0-2:",
			4,
			"invalid group size or used size: ",
		},
		{
			"0-2/8:10",
			8,
			"invalid end of range: 2/8",
		},
		{
			"a-2/8:10",
			8,
			"invalid start of range: a",
		},
		{
			"0-2a:10",
			8,
			"invalid end of range: 2a",
		},
		{
			"0-2:10",
			8,
			"invalid group size or used size: 10",
		},
		{
			"0-2:0/8",
			8,
			"used size must be at least 1, got: 0",
		},
		{
			"0-3:9/10",
			8,
			"used size greater than total CPUs: 9",
		},
	}

	for _, tt := range tst {
		t.Run(tt.input, func(t *testing.T) {
			_, err := i.ParseCPUs(tt.input, tt.tCores)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if err.Error() != tt.err {
				t.Fatalf("expected %v, got %v", tt.err, err)
			}
		})
	}
}

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
			res, err := i.GenerateComplementCPUList(tt.input, tt.tCores)
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

			check, err := i.AreCPUListsExclusiveWithMaxCPUs(tt.s1,
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
