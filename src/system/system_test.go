package system

import (
	"os"
	"testing"
)

var (
	osStat = os.Stat
	// osReadFile = os.ReadFile
)

// mockFile creates a temporary file with content and returns its path
// func mockFile(t *testing.T, path, content string) func() {
// 	t.Helper()

// 	// Create parent dirs if needed
// 	if err := os.MkdirAll(strings.TrimSuffix(path, "/model"), 0755); err != nil {
// 		t.Fatalf("failed to mkdir: %v", err)
// 	}

// 	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
// 		t.Fatalf("failed to write file: %v", err)
// 	}

// 	return func() {
// 		os.Remove(path)
// 	}
// }

func TestDetectSystemUbuntuCore(t *testing.T) {
	os.Setenv("SNAP_SAVE_DATA", "/some/uc/path")
	defer os.Unsetenv("SNAP_SAVE_DATA")

	sys, err := DetectSystem()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sys != UbuntuCore {
		t.Fatalf("expected UbuntuCore, got %v", sys)
	}
}

// func TestDetectSystem_Rpi(t *testing.T) {
// 	os.Unsetenv("SNAP_SAVE_DATA")

// 	cleanup := mockFile(t, "/tmp/device-tree/model", "Raspberry Pi 4 Model B")
// 	defer cleanup()

// 	// Temporarily override the real path
// 	originalStat := osStat
// 	originalRead := osReadFile
// 	defer func() {
// 		osStat = originalStat
// 		osReadFile = originalRead
// 	}()

// 	osStat = func(name string) (os.FileInfo, error) {
// 		if name == "/proc/device-tree/model" {
// 			return os.Stat("/tmp/device-tree/model")
// 		}
// 		return os.Stat(name)
// 	}

// 	osReadFile = func(name string) ([]byte, error) {
// 		if name == "/proc/device-tree/model" {
// 			return os.ReadFile("/tmp/device-tree/model")
// 		}
// 		return os.ReadFile(name)
// 	}

// 	sys, err := DetectSystem()
// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}
// 	if sys != Rpi {
// 		t.Fatalf("expected Rpi, got %v", sys)
// 	}
// }

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

	// Temporarily override os.Stat
	originalStat := osStat
	defer func() { osStat = originalStat }()

	osStat = func(name string) (os.FileInfo, error) {
		if name == "/etc/default/grub" {
			return os.Stat(tmpGrub)
		}
		return os.Stat(name)
	}

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
