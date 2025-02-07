package components

import (
	"log"
)

type IndexNav struct {
	current *int
	total   *int
}

func NewNavigation(current *int, total *int) *IndexNav {
	return &IndexNav{
		current: current,
		total:   total,
	}
}

// updateFocusIndex updates the focus index based on the given direction.
// direction should be +1 (move down) or -1 (move up).
// total is the total number of navigable items.
func (n *IndexNav) updateFocusIndex(direction int) {
	// Implement the circular navigation
	log.Println("\tCurrent: ", *n.current)

	*n.current = (*n.current + direction + *n.total) % *n.total
	log.Println("\t    New: ", *n.current)
}

func (n *IndexNav) Next() {
	n.updateFocusIndex(+1)
}

func (n *IndexNav) Prev() {
	n.updateFocusIndex(-1)
}
