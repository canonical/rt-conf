package model

import (
	"errors"
	"os"
	"testing"
)

func TestIRQTuningValidate(t *testing.T) {
	tests := []struct {
		name    string
		c       IRQTuning
		wantErr bool
	}{
		{
			name: "valid test",
			c: IRQTuning{
				CPUs: "0,1",
				Filter: IRQFilter{
					Actions:  `action`,
					ChipName: `chip-name`,
					Name:     `name`,
					Type:     `type`,
				},
			},
			wantErr: false,
		},
		{
			name: "valid regex",
			c: IRQTuning{
				CPUs: "0-N",
				Filter: IRQFilter{
					Actions:  `nvme`,
					ChipName: `-PCI-`,
					Name:     `name`,
					Type:     `type`,
				},
			},
			wantErr: false,
		},
		{
			name: "valid regex",
			c: IRQTuning{
				CPUs: "0,N",
				Filter: IRQFilter{
					Actions:  `nvme`,
					ChipName: `-PCI-`,
					Name:     `\d`,
					Type:     `type`,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid regex",
			c: IRQTuning{
				CPUs: "0,1",
				Filter: IRQFilter{
					Actions:  `(?!abc)def`, // Negative lookahead
					ChipName: `chip-name`,
					Name:     `name`,
					Type:     `type`,
				},
			},
			wantErr: true,
		},
	}

	t.Log("Running IRQTuning.Validate() tests")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Validate(); (err != nil) != tt.wantErr {
				t.Logf("Testing: %v\nFilters:\n\t%v", tt.name, tt.c.Filter)
				t.Errorf("IRQTuning.Validate() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

type mockDirEntry struct {
	name  string
	isDir bool
}

func (m *mockDirEntry) Name() string {
	return m.name
}

func (m *mockDirEntry) IsDir() bool {
	return m.isDir
}

func (m *mockDirEntry) Type() os.FileMode {
	return os.ModeDir
}

func (m *mockDirEntry) Info() (os.FileInfo, error) {
	return nil, nil
}

func TestGetHigherIRQ(t *testing.T) {
	tests := []struct {
		name       string
		entries    []os.DirEntry
		expected   int
		errReadDir error
		err        error
	}{
		{
			name:       "error reading dir",
			entries:    nil,
			errReadDir: errors.New("error reading dir"),
			err:        errors.New("error reading dir"),
		},
		{
			name: "invalid IRQ number",
			entries: []os.DirEntry{
				&mockDirEntry{name: "NonNumber", isDir: true},
				&mockDirEntry{name: "123", isDir: true},
			},
			err:      nil,
			expected: 123,
		},
		{
			name:       "no IRQs found",
			entries:    []os.DirEntry{},
			errReadDir: nil,
			err:        errors.New("no IRQs found"),
		},
		{
			name: "valid test",
			entries: []os.DirEntry{
				&mockDirEntry{name: "100", isDir: true},
				&mockDirEntry{name: "142", isDir: true},
				&mockDirEntry{name: "234", isDir: true},
			},
			errReadDir: nil,
			err:        nil,
			expected:   234,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			readDir = func(_ string) ([]os.DirEntry, error) {
				return tc.entries, tc.err
			}
			num, err := GetHigherIRQ()
			if tc.err != nil {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tc.err)
				}
				if err.Error() != tc.err.Error() {
					t.Fatalf("expected error %q, got %v", tc.err, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			if num != tc.expected {
				t.Fatalf("expected %d, got %d", tc.expected, num)
			}
		})
	}
}
