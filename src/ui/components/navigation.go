package components

import (
	"log"
	"sync"

	"github.com/canonical/rt-conf/src/ui/config"
)

// Singleton instance
var (
	instance *MenuNav
	once     sync.Once
)

type MenuNav struct {
	menuStack []config.Views
	// mut       sync.Mutex
}

// GetInstance ensures that only one instance of MenuNav is created
func GetMenuNavInstance() *MenuNav {
	once.Do(func() {
		log.Println("----> Initializing Singleton MenuNav")
		instance = &MenuNav{
			menuStack: []config.Views{config.INIT_VIEW_ID},
		}
	})
	return instance
}

func (m *MenuNav) PrintMenuStack() []config.Views {
	return m.menuStack
}

// Set a new menu and track the previous menu
func (m *MenuNav) SetNewMenu(newMenu config.Views) {
	log.Println("Current screen: ", config.Menu[m.GetCurrMenu()])
	log.Println("Going to screen: ", config.Menu[newMenu])
	m.menuStack = append(m.menuStack, newMenu)
}

// Get the current menu
func (m *MenuNav) GetCurrMenu() config.Views {
	return m.menuStack[len(m.menuStack)-1]
}

// Return to the previous menu
func (m *MenuNav) PrevMenu() {
	// m.mut.Lock()
	// defer m.mut.Unlock()

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
	// m.mut.Lock()
	// defer m.mut.Unlock()
	log.Println("----> ReturnToMainMenu Called")

	log.Println("Old menu: ", m.menuStack)
	m.menuStack = nil // Clear history
	m.menuStack = append(m.menuStack, config.INIT_VIEW_ID)
	log.Println("New menu: ", m.menuStack)
}
