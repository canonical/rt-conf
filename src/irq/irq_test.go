package irq

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/canonical/rt-conf/src/model"
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

func TestHappyIRQtuning(t *testing.T) {

	var happyCases = []IRQTestCase{
		{
			Yaml: `
irq_tuning:
- cpus: 0
  filter:
    action: floppy
`,
			Handler: &mockIRQReaderWriter{
				IRQs: map[uint]IRQInfo{
					0: {
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

func TestUnhappyIRQtuning(t *testing.T) {

	var UnhappyCases = []IRQTestCase{
		{
			// Invalid number
			Yaml: `
irq_tuning:
- cpus: 0
  filter:
    number: a
`,
			Handler: &mockIRQReaderWriter{},
		},
		{
			// Invalid RegEx
			Yaml: `
irq_tuning:
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
			if err == nil {
				t.Fatalf("Expected error, got nil on YAML %v", c.Yaml)
			}
		})
	}
}

func mainLogicIRQ(t *testing.T, cfg IRQTestCase, i int) (string, error) {
	tempConfigPath := setupTempFile(t, cfg.Yaml, i)
	t.Cleanup(func() {
		os.Remove(tempConfigPath)
	})
	var conf model.InternalConfig
	if d, err := model.LoadConfigFile(tempConfigPath); err != nil {
		return "", fmt.Errorf("failed to load config file: %v", err)
	} else {
		conf.Data = *d
	}

	err := applyIRQConfig(&conf, cfg.Handler)
	if err != nil {
		return "", fmt.Errorf("Failed to process interrupts: %v", err)
	}
	return "", nil
}

func TestWriteCPUAffinitySuccessfulWrite(t *testing.T) {
	tmpDir := t.TempDir()

	irqNum := 1
	cpus := "0-3"
	irqPath := filepath.Join(tmpDir, fmt.Sprintf("%d", irqNum))
	if err := os.MkdirAll(irqPath, 0755); err != nil {
		t.Fatalf("failed to create IRQ directory: %v", err)
	}
	affinityFile := filepath.Join(irqPath, "smp_affinity_list")
	f, err := os.Create(affinityFile)
	if err != nil {
		t.Fatalf("failed to create file: %v", err)
	}
	f.Close()

	procIRQ = tmpDir // override to avoid touching /proc
	writer := &realIRQReaderWriter{
		FileWriter: realFileWriter{},
	}
	err = writer.WriteCPUAffinity(irqNum, cpus)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(affinityFile)
	if err != nil {
		t.Fatalf("error reading back: %v", err)
	}
	if string(content) != cpus {
		t.Errorf("expected %q, got %q", cpus, string(content))
	}
}

type mockFileWriterIOError struct{}

func (mockFileWriterIOError) WriteFile(path string, content []byte, perm os.FileMode) error {
	return fmt.Errorf("input/output error") // Simulated /proc error
}

func TestWriteCPUAffinityInputOutputErrorIgnored(t *testing.T) {
	writer := &realIRQReaderWriter{
		FileWriter: mockFileWriterIOError{},
	}
	err := writer.WriteCPUAffinity(1, "0")
	if err != nil {
		t.Fatalf("expected nil, got error: %v", err)
	}
}

// Simulate a real write error that's not ignorable (not "input/output error")
func TestWriteCPUAffinityErrorReturned(t *testing.T) {
	procIRQ = "/this/path/does/not/exist"

	writer := &realIRQReaderWriter{
		FileWriter: realFileWriter{},
	}
	err := writer.WriteCPUAffinity(99, "1-2")

	if err == nil {
		t.Fatal("expected an error but got nil")
	}
	if !strings.Contains(err.Error(), "error writing to") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// Sanity: return nil even if file already has the value
func TestWriteCPUAffinityAlreadySet(t *testing.T) {
	tmpDir := t.TempDir()
	procIRQ = tmpDir

	irqNum := 5
	cpus := "0"
	irqPath := filepath.Join(tmpDir, fmt.Sprintf("%d", irqNum))
	if err := os.MkdirAll(irqPath, 0755); err != nil {
		t.Fatal(err)
	}
	affinityFile := filepath.Join(irqPath, "smp_affinity_list")
	if err := os.WriteFile(affinityFile, []byte(cpus), 0644); err != nil {
		t.Fatal(err)
	}

	writer := &realIRQReaderWriter{
		FileWriter: realFileWriter{},
	}
	err := writer.WriteCPUAffinity(irqNum, cpus)

	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

type irqDirEntry struct {
	Number int
	Files  map[string]string
}

func setupIRQTestDir(t *testing.T, entries []irqDirEntry) string {
	t.Helper()

	tmpDir := t.TempDir()
	sysKernelIRQ = tmpDir

	for _, e := range entries {
		dir := filepath.Join(tmpDir, strconv.Itoa(e.Number))
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}
		for name, content := range e.Files {
			if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
				t.Fatalf("failed to write file: %v", err)
			}
		}
	}
	return tmpDir
}
func TestReadIRQsSingleActiveIRQ(t *testing.T) {
	setupIRQTestDir(t, []irqDirEntry{
		{
			Number: 10,
			Files: map[string]string{
				"actions":   "handle_irq",
				"chip_name": "testchip",
				"name":      "eth0",
				"type":      "level",
				"wakeup":    "enabled",
			},
		},
	})

	r := &realIRQReaderWriter{}
	irqs, err := r.ReadIRQs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(irqs) != 1 {
		t.Fatalf("expected 1 irq, got %d", len(irqs))
	}
	if irqs[0].Number != 10 || irqs[0].Name != "eth0" {
		t.Fatalf("unexpected irq info: %+v", irqs[0])
	}
}

func TestReadIRQsEmptyActionsIgnored(t *testing.T) {
	setupIRQTestDir(t, []irqDirEntry{
		{
			Number: 11,
			Files: map[string]string{
				"actions": "",
			},
		},
	})

	r := &realIRQReaderWriter{
		FileWriter: realFileWriter{},
	}
	irqs, err := r.ReadIRQs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(irqs) != 0 {
		t.Fatalf("expected 0 irq, got %d", len(irqs))
	}
}

func TestReadIRQsNonNumericDirectoryIgnored(t *testing.T) {
	tmp := t.TempDir()
	sysKernelIRQ = tmp
	_ = os.Mkdir(filepath.Join(tmp, "notanumber"), 0755)

	r := &realIRQReaderWriter{
		FileWriter: realFileWriter{},
	}
	irqs, err := r.ReadIRQs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(irqs) != 0 {
		t.Fatalf("expected 0 irq, got %d", len(irqs))
	}
}

func TestReadIRQsReadDirError(t *testing.T) {
	sysKernelIRQ = "/invalid/path"

	r := &realIRQReaderWriter{
		FileWriter: realFileWriter{},
	}

	_, err := r.ReadIRQs()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestReadIRQsReadFileErrorHandled(t *testing.T) {
	setupIRQTestDir(t, []irqDirEntry{
		{
			Number: 12,
			Files: map[string]string{
				// Only one file, others will cause ReadFile errors
				"actions": "handle_irq",
			},
		},
	})

	r := &realIRQReaderWriter{
		FileWriter: realFileWriter{},
	}
	irqs, err := r.ReadIRQs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(irqs) != 1 {
		t.Fatalf("expected 1 irq, got %d", len(irqs))
	}
}

func TestApplyIRQConfig(t *testing.T) {
	config := &model.InternalConfig{
		Data: model.Config{
			Interrupts: []model.IRQTuning{},
		},
	}
	err := ApplyIRQConfig(config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
