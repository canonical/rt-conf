package pwrmgmt

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/canonical/rt-conf/src/cpulists"
	"github.com/canonical/rt-conf/src/model"
)

// setupTempDirWithFiles creates a temporary directory and then creates n files
// named "0", "1", ..., "n-1" inside that directory. It fails the test if any
// error occurs.
func setupTempDirWithFiles(t *testing.T, prvRule string, maxCpus int) string {
	t.Helper()

	tempDir, err := os.MkdirTemp("", "tempfiles-")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}

	// Clean up the temp directory after the test finishes.
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})

	// Create files from 0 to n-1.
	for i := 0; i < maxCpus; i++ {
		filename := strconv.Itoa(i)
		cpuPath := filepath.Join(tempDir, filename)
		if err := os.Mkdir(cpuPath, 0755); err != nil {
			t.Fatalf("failed to create directory %s: %v", cpuPath, err)
		}

		scalGov := filepath.Join(cpuPath, "scalgov")
		fscalGov, err := os.Create(scalGov)
		if err != nil {
			t.Fatalf("failed to create file %s: %v", scalGov, err)
		}

		nb, err := fscalGov.Write([]byte(prvRule))
		if err != nil {
			t.Fatalf("failed to write to file %s: %v", scalGov, err)
		}
		if nb != len(prvRule) {
			t.Fatalf("number of written bytes doesn't match on file %s",
				scalGov)
		}
		fscalGov.Close()

		for _, file := range []string{"maxfreq", "minfreq"} {
			filePath := filepath.Join(cpuPath, file)
			if err := os.WriteFile(filePath, []byte("0"), 0644); err != nil {
				t.Fatalf("failed to create file %s: %v", filePath, err)
			}
		}
	}

	return tempDir
}

func TestPwrMgmt(t *testing.T) {
	// Since this considers the real amount of cpus in the system, all cpulists
	// for CpuGovernanceRule.CPUs are set to 0 so it can be tested with any
	// amount of cpus
	var happyCases = []struct {
		name     string
		maxCpus  int
		prevRule string
		d        model.PwrMgmt // add only one rule here
	}{
		{
			"powersave to performance",
			3,
			"powersave",
			model.PwrMgmt{
				"foo": {
					CPUs:    "0",
					ScalGov: "performance",
				},
			},
		},
		{
			"performance to powersave",
			8,
			"performance",
			model.PwrMgmt{
				"bar": {
					CPUs:    "0",
					ScalGov: "balanced",
				},
			},
		},
		{
			"balanced to powersave",
			4,
			"balanced",
			model.PwrMgmt{
				"baz": {
					CPUs:    "0",
					ScalGov: "powersave",
				},
			},
		},
		{
			"balanced to powersave with min and max freq",
			4,
			"balanced",
			model.PwrMgmt{
				"buz": {
					CPUs:    "0",
					ScalGov: "powersave",
					MinFreq: "5.45GHz",
					MaxFreq: "5.584GHz",
				},
			},
		},
		{
			"balanced to powersave with min and max freq with different case",
			4,
			"balanced",
			model.PwrMgmt{
				"qux": {
					CPUs:    "0",
					ScalGov: "powersave",
					MinFreq: "2.1ghz",
					MaxFreq: "2.5GHZ",
				},
			},
		},
		{
			"balanced to powersave with only min freq",
			4,
			"balanced",
			model.PwrMgmt{
				"foobar": {
					CPUs:    "0",
					ScalGov: "powersave",
					MinFreq: "2.1ghz",
				},
			},
		},
		{
			"performance to powersave with only max freq",
			4,
			"performance",
			model.PwrMgmt{
				"quux": {
					CPUs:    "0",
					ScalGov: "powersave",
					MaxFreq: "4000mHz",
				},
			},
		},
		{
			"Only max freq set",
			4,
			"powersave",
			model.PwrMgmt{
				"corge": {
					CPUs:    "0",
					MaxFreq: "4gHz",
				},
			},
		},
		{
			"Only max and min freq set",
			4,
			"powersave",
			model.PwrMgmt{
				"grault": {
					CPUs:    "0",
					MaxFreq: "4gHz",
					MinFreq: "1gHz",
				},
			},
		},
	}

	for index, tc := range happyCases {
		t.Run(fmt.Sprintf("case-%d", index), func(t *testing.T) {

			basePath := setupTempDirWithFiles(t, tc.prevRule, tc.maxCpus)

			// Create a new ReaderWriter instance with the base path
			pwrmgmtReaderWriter.ScalingGovernorPath = basePath + "/%d/scalgov"
			pwrmgmtReaderWriter.MinFreqPath = basePath + "/%d/minfreq"
			pwrmgmtReaderWriter.MaxFreqPath = basePath + "/%d/maxfreq"

			err := pwrmgmtReaderWriter.applyPwrConfig(tc.d)
			if err != nil {
				t.Fatalf("error: %v", err)
			}

			for idx, rule := range tc.d {

				parsedCpus, err := cpulists.Parse(rule.CPUs)
				if err != nil {
					t.Fatalf("error parsing cpus: %v", err)
				}
				for cpu := range parsedCpus {
					content, err := os.ReadFile(
						filepath.Join(basePath, strconv.Itoa(cpu), "scalgov"))
					if err != nil {
						t.Fatalf("error reading file: %v", err)
					}
					if string(content) != tc.d[idx].ScalGov && tc.d[idx].ScalGov != "" {
						t.Fatalf("expected %q, got %q", tc.d[idx].ScalGov,
							string(content))
					}
				}

			}

		})
	}

	var UnhappyCases = []struct {
		name string
		cfg  *model.InternalConfig
		err  error
	}{
		{
			name: "Invalid CPU list",
			cfg: &model.InternalConfig{
				Data: model.Config{
					CpuGovernance: model.PwrMgmt{
						"foo": {
							CPUs:    "2-1",
							ScalGov: "performance",
						},
					},
				},
			},
			err: fmt.Errorf("start of range greater than end: 2-1"),
		},
	}
	for _, tc := range UnhappyCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ApplyPwrConfig(tc.cfg)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if err.Error() != tc.err.Error() {
				t.Fatalf("expected error: %q, got: %q", tc.err.Error(), err)
			}
		})
	}
}

func TestEmptyPwrMgmtRules(t *testing.T) {
	var errorCases = []struct {
		name string
		cfg  *model.InternalConfig
	}{
		{
			name: "No CPU Governance rules",
			cfg:  &model.InternalConfig{},
		},
	}

	for _, tc := range errorCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ApplyPwrConfig(tc.cfg)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func TestParseFreq(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  int
		expectErr string
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: -1,
		},
		{
			name:      "invalid raw kHz no unit",
			input:     "1000",
			expected:  1000,
			expectErr: "invalid format:",
		},
		{
			name:     "kHz uppercase",
			input:    "1000KHz",
			expected: 1000,
		},
		{
			name:     "valid raw Hz with suffix",
			input:    "100000Hz",
			expected: 100,
		},
		{
			name:      "invalid MHz lowercase",
			input:     "2.5m",
			expectErr: "invalid format:",
		},
		{
			name:     "MHz uppercase",
			input:    "2.5MHz",
			expected: 2500,
		},
		{
			name:      "invalid GHz lowercase",
			input:     "2.0g",
			expectErr: "invalid format:",
		},
		{
			name:     "GHz uppercase",
			input:    "2.0GHz",
			expected: 2_000_000,
		},
		{
			name:      "Invalid float",
			input:     "fooGHz",
			expectErr: "failed to parse frequency value:",
		},
		{
			name:      "Invalid format",
			input:     "123.4.5Mhz",
			expectErr: "failed to parse frequency value:",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := model.ParseFreq(tc.input)
			if tc.expectErr == "" {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if got != tc.expected {
					t.Errorf("expected %d, got %d", tc.expected, got)
				}
			} else {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tc.expectErr)
				}
				if !strings.Contains(err.Error(), tc.expectErr) {
					t.Errorf("expected error to contain %q, got %q", tc.expectErr, err.Error())
				}
			}
		})
	}
}

func TestWriteOnly(t *testing.T) {
	tmpDir := t.TempDir()

	testCases := []struct {
		name        string
		prepare     func(path string) // setup, e.g., permissions
		path        string
		data        string
		expectError bool
	}{
		{
			name: "success",
			prepare: func(path string) {
				// Create the file with write permission
				_ = os.WriteFile(path, []byte("old content"), 0644)
			},
			path:        filepath.Join(tmpDir, "success.txt"),
			data:        "new data",
			expectError: false,
		},
		{
			name: "fail to open (no such file)",
			prepare: func(path string) {
				// Do not create the file
			},
			path:        filepath.Join(tmpDir, "doesnotexist.txt"),
			data:        "won't matter",
			expectError: true,
		},
		{
			name: "fail to write (read-only file)",
			prepare: func(path string) {
				_ = os.WriteFile(path, []byte("content"), 0444) // Read-only
			},
			path:        filepath.Join(tmpDir, "readonly.txt"),
			data:        "",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(tc.path)

			err := writeOnly(tc.path, tc.data)
			if tc.expectError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("did not expect error, but got: %v", err)
			}
		})
	}
}

func TestApplyRuleUnhappy(t *testing.T) {
	tests := []struct {
		name   string
		sclgov model.CpuGovernanceRule
	}{
		{
			name: "WriteScalingGov fails (needs root)",
			sclgov: model.CpuGovernanceRule{
				ScalGov: "performance",
			},
		},
		{
			name: "ParseFreq MinFreq fails (invalid format)",
			sclgov: model.CpuGovernanceRule{
				MinFreq: "invalid!",
			},
		},
		{
			name: "ParseFreq MaxFreq fails (invalid format)",
			sclgov: model.CpuGovernanceRule{
				MinFreq: "1GHz",
				MaxFreq: "oops!",
			},
		},
		{
			name: "WriteCPUFreq fails on min frequency (needs root)",
			sclgov: model.CpuGovernanceRule{
				MinFreq: "1GHz",
			},
		},
		{
			name: "WriteCPUFreq fails on max frequency (needs root)",
			sclgov: model.CpuGovernanceRule{
				MaxFreq: "5GHz",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			err := pwrmgmtReaderWriter.applyRule(0, tc.sclgov)
			if err == nil {
				t.Fatalf("expected error got nil")
			}
		})
	}
}
