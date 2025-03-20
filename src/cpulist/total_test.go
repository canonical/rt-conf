package cpulist

import (
	"testing"
)

func TestTotalAvailable(t *testing.T) {
	c, err := TotalAvailable()
	if err != nil {
		t.Fatalf("Failed TotalAvailable: %v", err)
	}
	if c == 0 {
		t.Fatalf("Total CPUs is 0")
	}
}
