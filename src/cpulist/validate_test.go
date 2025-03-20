package cpulist

import "testing"

func TestHappyValidationWithFlags(t *testing.T) {
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
		{"0,n", max, isolcpuFlags},
		{"nohz,0,n", max, isolcpuFlags},
		{"domain,0,n", max, isolcpuFlags},
		{"managed_irq,0,n", max, isolcpuFlags},

		// Test comma separated CPU list
		{"0,n", max, isolcpuFlags},
		{"nohz,0,n", max, isolcpuFlags},
		{"domain,0,n", max, isolcpuFlags},
		{"managed_irq,0,n", max, isolcpuFlags},
	}

	for _, tc := range testCases {
		t.Run("TestValidationWithFlags", func(t *testing.T) {
			err := validateListWithFlags(tc.value, tc.flags, tc.cpus)
			if err != nil {
				t.Fatalf("Failed ValidateListWithFlags: %v", err)
			}
		})
	}
}
