package ui

import (
	"sync"
)

type menuItems struct {
	titles     []string
	descs      []string
	titleIndex int
	descIndex  int
	mtx        *sync.Mutex
}

func (r *menuItems) reset() {
	r.mtx = &sync.Mutex{}

	r.titles = []string{
		"Kernel cmdline",
		"IRQ Affinity",
	}

	r.descs = []string{
		"Configure Kernel cmdline parameters",
		"Isolate CPUs from serving IRQs",
	}
}

func (r *menuItems) next() item {
	if r.mtx == nil {
		r.reset()
	}

	r.mtx.Lock()
	defer r.mtx.Unlock()

	i := item{
		title:       r.titles[r.titleIndex],
		description: r.descs[r.descIndex],
	}

	r.titleIndex++
	if r.titleIndex >= len(r.titles) {
		r.titleIndex = 0
	}

	r.descIndex++
	if r.descIndex >= len(r.descs) {
		r.descIndex = 0
	}

	return i
}
