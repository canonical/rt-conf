package utils

import (
	"testing"
)

func TestTrimSurroundingQuotes(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"quoted"`, "quoted"},
		{`"unmatched`, `"unmatched`},
		{`no quotes`, `no quotes`},
		{`"spaces in here"`, `spaces in here`},
		{`"double quotes"`, `double quotes`},
	}

	for _, test := range tests {
		result := TrimSurroundingDoubleQuotes(test.input)
		if result != test.expected {
			t.Errorf("TrimSurroundingQuotes(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}
