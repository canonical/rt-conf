package irq

import (
	"os"
)

// IRQReaderWriter is an interface for read and write IRQ data from the filesystem.
type IRQReaderWriter interface {
	ReadIRQs() ([]IRQInfo, error)
	WriteCPUAffinity(irqNum int, cpus string) error
}

// Abstract the write operation
type FileWriter interface {
	WriteFile(path string, content []byte, perm os.FileMode) error
}

type realFileWriter struct{}

func (realFileWriter) WriteFile(path string, content []byte, perm os.FileMode) error {
	return os.WriteFile(path, content, perm)
}

// IRQInfo represents information about an IRQ.
type IRQInfo struct {
	// TODO: review these data types
	Number   int
	Actions  string
	ChipName string
	Name     string
	Type     string
	Wakeup   string
	// PerCPuCount string // ** NOTE: Not needed for now
}
