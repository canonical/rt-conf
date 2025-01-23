package cpu

import "testing"

func TestHappyValidationWithFlags(t *testing.T) {
	const max = 24
	var testCases = []struct {
		value string
		cpus  int
		flags []string
	}{
		// Test CPU list with range
		{"0-1", max, []string{"domain", "nohz", "managed_irq"}},
		{"nohz,0-1", max, []string{"domain", "nohz", "managed_irq"}},
		{"domain,0-1", max, []string{"domain", "nohz", "managed_irq"}},
		{"managed_irq,0-1", max, []string{"domain", "nohz", "managed_irq"}},

		// Test single CPU on CPU list
		{"0", max, []string{"domain", "nohz", "managed_irq"}},
		{"nohz,0", max, []string{"domain", "nohz", "managed_irq"}},
		{"domain,0", max, []string{"domain", "nohz", "managed_irq"}},
		{"managed_irq,0", max, []string{"domain", "nohz", "managed_irq"}},

		// Test comma separated CPU list
		{"0,n", max, []string{"domain", "nohz", "managed_irq"}},
		{"nohz,0,n", max, []string{"domain", "nohz", "managed_irq"}},
		{"domain,0,n", max, []string{"domain", "nohz", "managed_irq"}},
		{"managed_irq,0,n", max, []string{"domain", "nohz", "managed_irq"}},

		// Test comma separated CPU list
		{"0,n", max, []string{"domain", "nohz", "managed_irq"}},
		{"nohz,0,n", max, []string{"domain", "nohz", "managed_irq"}},
		{"domain,0,n", max, []string{"domain", "nohz", "managed_irq"}},
		{"managed_irq,0,n", max, []string{"domain", "nohz", "managed_irq"}},
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
