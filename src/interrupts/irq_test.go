package interrupts

import (
	"fmt"
	"os"
	"testing"

	"github.com/canonical/rt-conf/src/data"
)

// MockIRQReaderWriter is a mock implementation of IRQReaderWriter for testing.
type mockIRQReaderWriter struct {
	IRQs            map[uint]IRQInfo
	WrittenAffinity map[int]string
	Errors          map[string]error
}

func (m *mockIRQReaderWriter) ReadIRQs() ([]IRQInfo, error) {
	if err, ok := m.Errors["ReadIRQs"]; ok {
		return nil, err
	}
	irqInfos := make([]IRQInfo, 0, len(m.IRQs))
	for _, info := range m.IRQs {
		irqInfos = append(irqInfos, info)
	}
	return irqInfos, nil
}

func (m *mockIRQReaderWriter) WriteCPUAffinity(irqNum int, cpus string) error {
	if err, ok := m.Errors["WriteCPUAffinity"]; ok {
		return err
	}
	if m.WrittenAffinity == nil {
		m.WrittenAffinity = make(map[int]string)
	}
	// TODO: Find a way to expose this to the test
	fmt.Printf("Writing affinity for IRQ %d: %s", irqNum, cpus)
	m.WrittenAffinity[irqNum] = cpus
	return nil
}

type IRQTestCase struct {
	Yaml    string
	Handler IRQReaderWriter
}

func setupTempFile(t *testing.T, content string, idex int) string {
	t.Helper()

	tmpFile, err := os.CreateTemp("", fmt.Sprintf("tempfile-%d", idex))
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	return tmpFile.Name()
}

func TestHappyIRQtunning(t *testing.T) {

	var happyCases = []IRQTestCase{
		{
			Yaml: `
irq_tunning:
- cpus: 0
  filter:
    action: floppy
`,
			Handler: &mockIRQReaderWriter{
				IRQs: map[uint]IRQInfo{
					10: {
						Actions: "floppy",
					},
				},
			},
		},
	}

	for i, c := range happyCases {
		t.Run("Happy Cases", func(t *testing.T) {
			_, err := mainLogicIRQ(t, c, i)
			if err != nil {
				t.Fatalf("On YAML: \n%v\nError: %v", c.Yaml, err)
			}
		})
	}
}

func TestUnhappyIRQtunning(t *testing.T) {

	var UnhappyCases = []IRQTestCase{
		{
			// Invalid number
			Yaml: `
irq_tunning:
- cpus: 0
  filter:
    number: a
`,
			Handler: &mockIRQReaderWriter{},
		},
		{
			// Invalid RegEx
			Yaml: `
irq_tunning:
- cpus: 0
  filter:
    number: 0
    action: "*"
`,
			Handler: &mockIRQReaderWriter{},
		},
	}

	for i, c := range UnhappyCases {
		t.Run("Unhappy Cases", func(t *testing.T) {
			_, err := mainLogicIRQ(t, c, i)
			// if err != nil {
			// 	t.Fatalf("On YAML: \n%v\nError: %v", c.Yaml, err)
			// }
			if err == nil {
				t.Fatalf("Expected error, got nil on YAML %v", c.Yaml)
			}
		})
	}
}

func mainLogicIRQ(t *testing.T, c IRQTestCase, i int) (string, error) {
	tempConfigPath := setupTempFile(t, c.Yaml, i)
	t.Cleanup(func() {
		os.Remove(tempConfigPath)
	})
	var conf data.InternalConfig
	if d, err := data.LoadConfigFile(tempConfigPath); err != nil {
		return "", fmt.Errorf("failed to load config file: %v", err)
	} else {
		conf.Data = *d
	}

	err := applyIRQConfig(&conf, c.Handler)
	if err != nil {
		return "", fmt.Errorf("Failed to process interrupts: %v", err)
	}
	return "", nil
}
