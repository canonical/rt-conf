package interrupts

// IRQReader is an interface for reading IRQ data from the filesystem.
type IRQReader interface {
	ReadIRQs() ([]IRQInfo, error)
}

// IRQWriter is an interface for writing IRQ affinity to the filesystem.
type IRQWriter interface {
	WriteCPUAffinity(irqNum, cpus string) error
}

// IRQInfo represents information about an IRQ.
type IRQInfo struct {
	// TODO: review these data types
	Number   uint
	Actions  string
	ChipName string
	Name     string
	Type     string
	Wakeup   string
	// PerCPuCount string // ** NOTE: Not needed for now
}
