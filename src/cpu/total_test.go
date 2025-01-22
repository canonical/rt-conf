package cpu_test

import (
	"testing"

	"github.com/canonical/rt-conf/src/cpu"
)

func TestTotalAvailable(t *testing.T) {
	c, err := cpu.TotalAvailable()
	if err != nil {
		t.Fatalf("Failed TotalAvailable: %v", err)
	}
	if c == 0 {
		t.Fatalf("Total CPUs is 0")
	}
}
