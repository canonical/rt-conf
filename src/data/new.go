package data

import (
	"log"

	"github.com/canonical/rt-conf/src/cpu"
)

func NewInternCfg() *InternalConfig {

	c, err := cpu.TotalAvailable()
	if err != nil {
		log.Fatalf("Failed to get total available CPUs: %v", err)
	}
	return &InternalConfig{
		TotalCPUs: c,
	}
}
