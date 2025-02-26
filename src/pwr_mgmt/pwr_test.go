package pwrmgmt

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/canonical/rt-conf/src/cpu"
	"github.com/canonical/rt-conf/src/data"
)

type mockScalGovReaderWriter struct {
	basePath string
}

func (m *mockScalGovReaderWriter) WriteScalingGov(sclgov string, cpu int) error {
	// Create the target file path: basePath/<cpu>
	scalingGovFile := filepath.Join(m.basePath, strconv.Itoa(cpu))

	err := os.WriteFile(scalingGovFile, []byte(sclgov), 0644)
	if err != nil {
		return fmt.Errorf("error writing to %s: %v", scalingGovFile, err)
	} else {
		log.Printf("Set %s to %s", scalingGovFile, sclgov)
	}
	return nil
}

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
		filePath := filepath.Join(tempDir, filename)
		f, err := os.Create(filePath)
		if err != nil {
			t.Fatalf("failed to create file %s: %v", filePath, err)
		}
		nb, err := f.Write([]byte(prvRule))
		if err != nil {
			t.Fatalf("failed to write to file %s: %v", filePath, err)
		}
		if nb != len(prvRule) {
			t.Fatalf("number of written bytes doesn't match on file %s",
				filePath)
		}
		f.Close()
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
		d        []data.CpuGovernanceRule // add only one rule here
	}{
		{3,
			"powersave",
			[]data.CpuGovernanceRule{
				{
					CPUs:    "0",
					ScalGov: "performance",
				},
			},
		},
		{
			8,
			"performance",
			[]data.CpuGovernanceRule{
				{
					CPUs:    "0",
					ScalGov: "balanced",
				},
			},
		},
		{
			4,
			"balanced",
			[]data.CpuGovernanceRule{
				{
					CPUs:    "0",
					ScalGov: "powersave",
				},
			},
		},
	}

	for index, tc := range happyCases {
		t.Run(fmt.Sprintf("case-%d", index), func(t *testing.T) {

			basePath := setupTempDirWithFiles(t, tc.prevRule, tc.maxCpus)
			scalingGovernerReaderWriter.Path = basePath + "/%d"
			err := scalingGovernerReaderWriter.applyPwrConfig(tc.d)
			if err != nil {
				t.Fatalf("error: %v", err)
			}

			for idx, rule := range tc.d {

				parsedCpus, err := cpu.ParseCPUs(rule.CPUs)
				if err != nil {
					t.Fatalf("error parsing cpus: %v", err)
				}
				for cpu := range parsedCpus {
					content, err := os.ReadFile(
						filepath.Join(basePath, strconv.Itoa(cpu)))
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
}
