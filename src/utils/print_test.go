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

func assertContains(t *testing.T, haystack, needle string) {
	t.Helper()
	if !strings.Contains(haystack, needle) {
		t.Errorf("expected to find %q in output: %q", needle, haystack)
	}
}

func TestPrintBoldBgText(t *testing.T) {
	output := captureLogOutput(func() {
		printBoldBgText("Hello %s", "world")
	})
	expectedSubstring := initColor + "Hello world" + endColor
	assertContains(t, output, expectedSubstring)
}

func TestPrintlnBoldBgText(t *testing.T) {
	output := captureLogOutput(func() {
		printlnBoldBgText("Test %d", 123)
	})
	expectedSubstring := initColor + "Test 123" + endColor
	assertContains(t, output, "\n")
	assertContains(t, output, expectedSubstring)
}

func TestPrintTitle(t *testing.T) {
	output := captureLogOutput(func() {
		PrintTitle("Header")
	})

	// We expect box-drawing characters and colored title line
	assertContains(t, output, "│ Header │")
	assertContains(t, output, initColor)
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
		assertContains(t, output, line)
	}
}
