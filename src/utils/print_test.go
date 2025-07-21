package utils

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

// Helper: capture stdout
func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
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
	output := captureStdout(func() {
		printBoldBgText("Hello %s", "world")
	})
	expectedSubstring := "\033[1;48;2;233;84;32mHello world\033[0m"
	if !strings.Contains(output, expectedSubstring) {
		t.Errorf("expected escape sequence with colored text, got: %q", output)
	}
}

func TestPrintlnBoldBgText(t *testing.T) {
	output := captureStdout(func() {
		printlnBoldBgText("Test %d", 123)
	})
	expectedSubstring := "\033[1;48;2;233;84;32mTest 123\033[0m"
	if !strings.Contains(output, expectedSubstring) {
		t.Errorf("expected colored output with newline, got: %q", output)
	}
}

func TestPrintTitle(t *testing.T) {
	output := captureStdout(func() {
		PrintTitle("Header")
	})
	// We expect box-drawing characters and colored title line
	if !strings.Contains(output, "│ Header │") {
		t.Errorf("title box missing expected content, got: %q", output)
	}
	if !strings.Contains(output, "\033[1;48;2;233;84;32m") {
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
