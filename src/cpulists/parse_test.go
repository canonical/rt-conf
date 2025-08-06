package cpulists

import (
	"errors"
	"reflect"
	"testing"
)

func TestParseCPUListsHappy(t *testing.T) {
	type test struct {
		input  string
		tCores int
		output CPUs
	}

	var tst = []test{
		// single CPU
		{
			"all",
			4,
			CPUs{0: true, 1: true, 2: true, 3: true},
		},
		{
			"all",
			2,
			CPUs{0: true, 1: true},
		},
		{
			"0-N",
			2,
			CPUs{0: true, 1: true},
		},
		{
			"4",
			8,
			CPUs{4: true},
		},
		// 3 single CPUs
		{
			"4,5,9",
			10,
			CPUs{4: true, 5: true, 9: true},
		},
		// CPU range
		{
			"0-7",
			10,
			CPUs{0: true, 1: true, 2: true, 3: true, 4: true,
				5: true, 6: true, 7: true},
		},
		// two CPUs ranges
		{
			"0-2,4-7",
			10,
			CPUs{0: true, 1: true, 2: true, 4: true,
				5: true, 6: true, 7: true},
		},
		// CPU range + single CPU
		{
			"0-2,3",
			4,
			CPUs{0: true, 1: true, 2: true, 3: true},
		},
		// Formated CPU list
		{
			"0-20:2/5",
			24,
			CPUs{0: true, 1: true, 5: true, 6: true,
				10: true, 11: true, 15: true, 16: true, 20: true},
		},
		// Formated CPU list + a single CPU
		{
			"0-20:2/5,23",
			24,
			CPUs{0: true, 1: true, 5: true, 6: true, 10: true,
				11: true, 15: true, 16: true, 20: true, 23: true},
		},
	}

	// If not set totalCPUs back to the original function, the next tests will fail.
	// Because totalCPUs is a global function pointer in this module.
	originalTotalCPUs := totalCPUs
	t.Cleanup(func() { totalCPUs = originalTotalCPUs })

	for _, tt := range tst {
		t.Run(tt.input, func(t *testing.T) {
			totalCPUs = func() (int, error) {
				return tt.tCores, nil
			}

			res, err := Parse(tt.input)
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

func TestParseCPUListsUnhappy(t *testing.T) {
	type test struct {
		input  string
		tCores int
		err    string
	}

	var tst = []test{
		{
			"al",
			2,
			"invalid CPU: al",
		},
		{
			"alll",
			2,
			"invalid CPU: alll",
		},
		{
			"0-all",
			4,
			"invalid end of range: all",
		},
		{
			"all-N",
			4,
			"invalid start of range: all",
		},
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
		{
			"0-n",
			8,
			"lowercase N isn't accepted",
		},
	}

	// If not set totalCPUs back to the original function, the next tests will fail.
	// Because totalCPUs is a global function pointer in this module.
	originalTotalCPUs := totalCPUs
	t.Cleanup(func() { totalCPUs = originalTotalCPUs })
	t.Run("TestParseCPUListsUnhappy", func(t *testing.T) {
		totalCPUs = func() (int, error) {
			return 0, errors.New("access denied")
		}
		expectedErr := "failed to get total available CPUs: access denied"
		_, err := Parse("0-1")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if err.Error() != expectedErr {
			t.Fatalf("expected '%v', got '%v'", expectedErr, err)
		}

	})

	for _, tt := range tst {
		t.Run(tt.input, func(t *testing.T) {
			_, err := ParseForCPUs(tt.input, tt.tCores)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if err.Error() != tt.err {
				t.Fatalf("expected %v, got %v", tt.err, err)
			}
		})
	}
}

func TestParseWithFlagsHappy(t *testing.T) {
	isolcpuFlags := []string{"domain", "nohz", "managed_irq"}
	const max = 24
	var testCases = []struct {
		value string
		cpus  int
		flags []string
	}{
		// Test CPU list with range
		{"0-1", max, isolcpuFlags},
		{"nohz,0-1", max, isolcpuFlags},
		{"domain,0-1", max, isolcpuFlags},
		{"managed_irq,0-1", max, isolcpuFlags},

		// Test single CPU on CPU list
		{"0", max, isolcpuFlags},
		{"nohz,0", max, isolcpuFlags},
		{"domain,0", max, isolcpuFlags},
		{"managed_irq,0", max, isolcpuFlags},

		// Test comma separated CPU list
		{"0,N", max, isolcpuFlags},
		{"nohz,0,N", max, isolcpuFlags},
		{"domain,0,N", max, isolcpuFlags},
		{"managed_irq,0,N", max, isolcpuFlags},

		// Test comma separated CPU list
		{"0,N", max, isolcpuFlags},
		{"nohz,0,N", max, isolcpuFlags},
		{"domain,0,N", max, isolcpuFlags},
		{"managed_irq,0,N", max, isolcpuFlags},
	}

	// If not set totalCPUs back to the original function, the next tests will fail.
	// Because totalCPUs is a global function pointer in this module.
	originalTotalCPUs := totalCPUs
	t.Cleanup(func() { totalCPUs = originalTotalCPUs })

	for _, tc := range testCases {
		t.Run("TestValidationWithFlags", func(t *testing.T) {
			totalCPUs = func() (int, error) {
				return tc.cpus, nil
			}
			_, _, err := ParseWithFlags(tc.value, tc.flags)
			if err != nil {
				t.Fatalf("Failed ValidateListWithFlags: %v", err)
			}
		})
	}

	t.Run("TestUnhappy-Failed-totalCPUs", func(t *testing.T) {
		totalCPUs = func() (int, error) {
			return 0, errors.New("access denied")
		}
		expectedErr := "failed to get total available CPUs: access denied"
		_, _, err := ParseWithFlags("", isolcpuFlags)
		if err == nil {
			t.Fatalf("Expected error: %v, got nil", expectedErr)
		}
		if err.Error() != expectedErr {
			t.Fatalf("Expected error: '%v', got '%v'", expectedErr, err)
		}
	})
}

func TestGenCPUlist(t *testing.T) {
	var testCases = []struct {
		name   string
		cpus   []int
		result string
	}{
		{
			name:   "TestEmptyCPUs",
			cpus:   []int{},
			result: "",
		},
		{
			name:   "TestSingleCPU",
			cpus:   []int{0},
			result: "0",
		},
		{
			name:   "TestMultipleSingleCPUs",
			cpus:   []int{0, 2, 4, 6},
			result: "0,2,4,6",
		},
		{
			name:   "TestCPURange",
			cpus:   []int{0, 1, 2},
			result: "0-2",
		},
		{
			name:   "TestMultipleCPURanges",
			cpus:   []int{0, 1, 2, 4, 5},
			result: "0-2,4-5",
		},
		{
			name:   "TestNonContiguousCPUs",
			cpus:   []int{0, 2, 4},
			result: "0,2,4",
		},
		{
			name:   "TestMixedCPUsRangeAndSingle",
			cpus:   []int{9, 1, 2, 3, 6},
			result: "1-3,6,9",
		},
		{
			name:   "TestMixedCPUsSingleAndRange",
			cpus:   []int{0, 2, 4, 6, 7},
			result: "0,2,4,6-7",
		},
		{
			name:   "TestCPUsWithGaps",
			cpus:   []int{0, 1, 3, 4, 6},
			result: "0-1,3-4,6",
		},
		{
			name:   "TestCPUsWithGapsAndRange",
			cpus:   []int{0, 1, 3, 4, 6, 7},
			result: "0-1,3-4,6-7",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GenCPUlist(tc.cpus)
			if result != tc.result {
				t.Errorf("Expected %s, got %s", tc.result, result)
			}
		})
	}
}
