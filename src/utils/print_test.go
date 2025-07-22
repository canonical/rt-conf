package utils

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	log.SetFlags(0)
	os.Exit(m.Run())
}

// Helper: capture log output
func captureLogOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	f()
	return buf.String()
}

func TestPrintBoldBgText(t *testing.T) {
	output := captureLogOutput(func() {
		printBoldBgText("Hello %s", "world")
	})
	expectedSubstring := initColor + "Hello world" + endColor
	if !strings.Contains(output, expectedSubstring) {
		t.Errorf("expected escape sequence with colored text, got: %q", output)
	}
}

func TestPrintlnBoldBgText(t *testing.T) {
	output := captureLogOutput(func() {
		printlnBoldBgText("Test %d", 123)
	})
	expectedSubstring := initColor + "Test 123" + endColor
	if !strings.Contains(output, expectedSubstring) {
		t.Errorf("expected: %q, got: %q", expectedSubstring, output)
	}
}

func TestPrintTitle(t *testing.T) {
	output := captureLogOutput(func() {
		PrintTitle("Header")
	})
	// We expect box-drawing characters and colored title line
	if !strings.Contains(output, "│ Header │") {
		t.Errorf("title box missing expected content, got: %q", output)
	}
	if !strings.Contains(output, initColor) {
		t.Errorf("expected background color escape sequence, got: %q", output)
	}
}

func TestLogTreeStyle(t *testing.T) {
	output := captureLogOutput(func() {
		LogTreeStyle([]string{"one", "two", "three"})
	})

	expected := []string{
		"├── one",
		"├── two",
		"└── three",
	}
	for _, line := range expected {
		if !strings.Contains(output, line) {
			t.Errorf("missing expected tree log entry: %q", line)
		}
	}
}
