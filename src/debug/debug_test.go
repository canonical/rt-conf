package debug

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestDebugPrintfAndPrintln(t *testing.T) {
	// nothing should be printed if debug is off
	buf := &bytes.Buffer{}
	log.SetOutput(buf)

	Printf("foo %d", 1)
	Println("bar")

	// check that nothing was printed
	if buf.Len() != 0 {
		t.Fatal("debug output not empty")
	}

	// enable and test that it actually logs
	Enable()
	Printf("hello %s", "world")
	Println("!")
	out := buf.String()
	if !strings.Contains(out, "hello world") || !strings.Contains(out, "!") {
		t.Fatal("debug output missing")
	}
}
