package interrupts

import "fmt"

// MockIRQReader is a mock implementation of IRQReader for testing.
type MockIRQReader struct {
	IRQs   map[uint]IRQInfo
	Errors map[string]error
}

func (m *MockIRQReader) ReadIRQs() ([]IRQInfo, error) {
	if err, ok := m.Errors["ReadIRQs"]; ok {
		return nil, err
	}
	irqInfos := make([]IRQInfo, 0, len(m.IRQs))
	for _, info := range m.IRQs {
		irqInfos = append(irqInfos, info)
	}
	return irqInfos, nil
}

// MockIRQWriter is a mock implementation of IRQWriter for testing.
type MockIRQWriter struct {
	WrittenAffinity map[int]string
	Errors          map[string]error
}

func (m *MockIRQWriter) WriteCPUAffinity(irqNum int, cpus string) error {
	fmt.Println("[DEBUG] MOCK WriteCPUAffinity")
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
