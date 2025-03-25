package cpulists

import (
	"testing"
)

func TestTotalCPUs(t *testing.T) {
	c, err := totalCPUs()
	if err != nil {
		t.Fatalf("Failed to get total CPUs: %s", err)
	}
	if c == 0 {
		t.Fatal("Unexpected total CPUs: 0")
	}
}
