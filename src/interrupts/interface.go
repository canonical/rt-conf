package interrupts

// IRQReaderWriter is an interface for read and write IRQ data from the filesystem.
type IRQReaderWriter interface {
	ReadIRQs() ([]IRQInfo, error)
	WriteCPUAffinity(irqNum int, cpus string) error
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
