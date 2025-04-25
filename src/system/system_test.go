package system

import (
	"os"
	"testing"
)

func TestDetectSystemUbuntuCore(t *testing.T) {
	os.Setenv("SNAP_SAVE_DATA", "/some/uc/path")
	t.Cleanup(func() {
		os.Unsetenv("SNAP_SAVE_DATA")
	})

	sys, err := DetectSystem()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sys != UbuntuCore {
		t.Fatalf("expected UbuntuCore, got %v", sys)
	}
}

func TestDetectSystemGrub(t *testing.T) {
	os.Unsetenv("SNAP_SAVE_DATA")

	tmpGrub := "/tmp/default/grub"
	if err := os.MkdirAll("/tmp/default", 0755); err != nil {
		t.Fatalf("failed to mkdir: %v", err)
	}
	if err := os.WriteFile(tmpGrub, []byte("GRUB_CMDLINE"), 0644); err != nil {
		t.Fatalf("failed to write grub file: %v", err)
	}
	defer os.Remove(tmpGrub)

	sys, err := DetectSystem()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sys != Grub {
		t.Fatalf("expected Grub, got %v", sys)
	}
}

func TestDetectSystemUnknown(t *testing.T) {
	baseDir = "/tmp"

	sys, err := DetectSystem()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sys != Unknown {
		t.Fatalf("expected Unknown, got %v", sys)
	}
}
