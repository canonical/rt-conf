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
		maxCpus  int
		prevRule string
		d        []model.CpuGovernanceRule // add only one rule here
	}{
		{3,
			"powersave",
			[]model.CpuGovernanceRule{
				{
					CPUs:    "0",
					ScalGov: "performance",
				},
			},
		},
		{
			8,
			"performance",
			[]model.CpuGovernanceRule{
				{
					CPUs:    "0",
					ScalGov: "balanced",
				},
			},
		},
		{
			4,
			"balanced",
			[]model.CpuGovernanceRule{
				{
					CPUs:    "0",
					ScalGov: "powersave",
				},
			},
		},
		{
			4,
			"balanced",
			[]model.CpuGovernanceRule{
				{
					CPUs:    "0",
					ScalGov: "powersave",
					MaxFreq: "5584000",
					MinFreq: "545000",
				},
			},
		},
		{
			4,
			"balanced",
			[]model.CpuGovernanceRule{
				{
					CPUs:    "0",
					ScalGov: "powersave",
					MaxFreq: "2.5GHz",
					MinFreq: "2.1GHz",
				},
			},
		},
		{
			4,
			"balanced",
			[]model.CpuGovernanceRule{
				{
					CPUs:    "0",
					ScalGov: "powersave",
					MaxFreq: "2.5G",
					MinFreq: "2.1G",
				},
			},
		},
		{
			4,
			"performance",
			[]model.CpuGovernanceRule{
				{
					CPUs:    "0",
					ScalGov: "powersave",
					MaxFreq: "4000mHz",
					MinFreq: "2000mHz",
				},
			},
		},
	}

	for index, tc := range happyCases {
		t.Run(fmt.Sprintf("case-%d", index), func(t *testing.T) {

			basePath := setupTempDirWithFiles(t, tc.prevRule, tc.maxCpus)

			// Create a new ReaderWriter instance with the base path
			pwrmgmtReaderWriter.ScalingGovernorPath = basePath + "/%d/scalgov"
			pwrmgmtReaderWriter.MaxFreqPath = basePath + "/%d/maxfreq"
			pwrmgmtReaderWriter.MinFreqPath = basePath + "/%d/minfreq"

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
					if string(content) != tc.d[idx].ScalGov {
						t.Fatalf("expected %s, got %s", tc.d[idx].ScalGov,
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
					CpuGovernance: []model.CpuGovernanceRule{
						{
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
				t.Fatalf("expected error: %v, got: %v", tc.err.Error(), err)
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

func TestCheckFequencyRules(t *testing.T) {
	tests := []struct {
		name      string
		min       int
		max       int
		expectErr string
	}{
		{
			name:      "no limits set",
			min:       0,
			max:       0,
			expectErr: "",
		},
		{
			name:      "valid frequency range",
			min:       1000,
			max:       2000,
			expectErr: "",
		},
		{
			name:      "negative max",
			min:       1000,
			max:       -2000,
			expectErr: "frequency values must be non-negative",
		},
		{
			name:      "negative min",
			min:       -1000,
			max:       2000,
			expectErr: "frequency values must be non-negative",
		},
		{
			name:      "max < min",
			min:       3000,
			max:       2000,
			expectErr: "max frequency (2000) cannot be less than min frequency (3000)",
		},
		{
			name:      "min == max",
			min:       1500,
			max:       1500,
			expectErr: "min and max frequency cannot be the same",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := CheckFequencyRules(tc.min, tc.max)
			if tc.expectErr == "" {
				if err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
			} else {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tc.expectErr)
				}
				if !strings.Contains(err.Error(), tc.expectErr) {
					t.Fatalf("expected error to contain %q, got: %v", tc.expectErr, err)
				}
			}
		})
	}
}

