package main

import (
	"fmt"
	"log"

	"github.com/canonical/rt-conf/src/interrupts"
)

func main() {

	irq, err := interrupts.MapIRQs()
	if err != nil {
		log.Fatal(err)
	}

	for _, irq := range irq {
		fmt.Printf("\nIRQ %d: %s\n", irq.Number, irq.Name)
		fmt.Printf("  CPUs: %v\n", irq.CPUs)
	}

	fmt.Printf("Size of IRQs: %v\n", len(irq))

	fmt.Println("\n\n\n\nSystem IRQs:")

	sysirqs, err := interrupts.MapSystemIRQs()
	if err != nil {
		log.Fatal(err)
	}
	for _, irq := range sysirqs {
		fmt.Printf("\nIRQ %d: %s\n", irq.Number, irq.Affinity)
	}

	fmt.Printf("Size of System IRQs: %v\n", len(sysirqs))
}
