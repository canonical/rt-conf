package ui

type menuItems struct {
	titles     []string
	descs      []string
	titleIndex int
	descIndex  int
}

func (r *menuItems) Size() int {
	return len(r.titles)
}

func (r *menuItems) Init() {

	r.titles = []string{
		"Kernel cmdline",
		"IRQ Affinity",
		"Power Management",
	}

	r.descs = []string{
		"Configure Kernel cmdline parameters",
		"Isolate CPUs from serving IRQs",
		"Configure CPU power management settings",
	}
}

func (r *menuItems) next() item {

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
