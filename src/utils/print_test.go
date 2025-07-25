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
	expectedSubstring := "\033[1mHello world\033[0m"
	assertContains(t, output, expectedSubstring)
}

func TestPrintlnBoldBgText(t *testing.T) {
	output := captureLogOutput(func() {
		printlnBoldBgText("Test %d", 123)
	})
	expectedSubstring := "\033[1mTest 123\033[0m"
	assertContains(t, output, "\n")
	assertContains(t, output, expectedSubstring)
}

func TestPrintTitle(t *testing.T) {
	output := captureLogOutput(func() {
		PrintTitle("Header")
	})

	assertContains(t, output, "│ Header │")
	assertContains(t, output, startBold)
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
