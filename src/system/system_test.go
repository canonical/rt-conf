package system

import (
	"os"
	"testing"
)

func TestDetectSystemUbuntuCore(t *testing.T) {
	if err := os.Setenv("SNAP_SAVE_DATA", "/some/uc/path"); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Unsetenv("SNAP_SAVE_DATA"); err != nil {
			t.Fatalf("failed to unset env: %v", err)
		}
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
	if err := os.Unsetenv("SNAP_SAVE_DATA"); err != nil {
		t.Fatalf("failed to unset env: %v", err)
	}

	tmpGrub := "/tmp/default/grub"
	if err := os.MkdirAll("/tmp/default", 0o755); err != nil {
		t.Fatalf("failed to mkdir: %v", err)
	}
	if err := os.WriteFile(tmpGrub, []byte("GRUB_CMDLINE"), 0o644); err != nil {
		t.Fatalf("failed to write grub file: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Remove(tmpGrub); err != nil {
			t.Fatal("failed to remove tmp grub file:", err)
		}
	})

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
