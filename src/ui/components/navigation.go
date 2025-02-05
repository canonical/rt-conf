package components

import (
	"log"

	"github.com/canonical/rt-conf/src/ui/config"
)

type MenuNav struct {
	menuStack []config.Views
	// This may be needed if the problem of the menu navigation is not fixed
	// mut       sync.Mutex
}

func NewMenuNav() *MenuNav {
	return &MenuNav{
		menuStack: []config.Views{
			config.INIT_VIEW_ID,
		},
	}
}

func (m *MenuNav) PrintMenuStack() []config.Views {
	return m.menuStack
}

// Set a new menu and track the previous menu
func (m *MenuNav) SetNewMenu(newMenu config.Views) {

	log.Println("----> SetNewMenu Called")
	log.Println("OLD: ", m.menuStack)

	m.menuStack = append(m.menuStack, newMenu)

	log.Println("NEW: ", m.menuStack)
}

// Get the current menu
func (m *MenuNav) GetCurrMenu() config.Views {
	return m.menuStack[len(m.menuStack)-1]
}

// Return to the previous menu
func (m *MenuNav) PrevMenu() {

	log.Println("----> PrevMenu Called")

	// it needs to have at least one menu which is the main menu
	if len(m.menuStack) > 1 {
		// Pop last menu from stack
		log.Println("Old menu: ", m.menuStack)
		m.menuStack = m.menuStack[:len(m.menuStack)-1]
		log.Println("New menu: ", m.menuStack)
	} else {
		log.Println("Old menu: ", m.menuStack)
		m.menuStack = nil
		m.menuStack = append(m.menuStack, config.INIT_VIEW_ID)
		log.Println("New menu: ", m.menuStack)
	}
}

// Reset back to mainMenu and clear history
func (m *MenuNav) ReturnToMainMenu() {
	log.Println("----> ReturnToMainMenu Called")

	log.Println("Old menu: ", m.menuStack)
	m.menuStack = nil // Clear history
	m.menuStack = append(m.menuStack, config.INIT_VIEW_ID)
	log.Println("New menu: ", m.menuStack)
}
