package pwrmgmt

import "github.com/canonical/rt-conf/src/data"

// IRQReaderWriter is an interface for read and write IRQ data from the filesystem.
type ScalGovReaderWriter interface {
	ReadPwrSetting() ([]PwrInfo, error) //TODO: maybe drop this
	WriteScalingGov(sclgov string, cpu int) error
}

// IRQInfo represents information about an IRQ.
type PwrInfo struct {
	Number  int
	ScalGov data.ScalProfiles
}
