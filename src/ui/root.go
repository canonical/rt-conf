package ui

import (
	"log"
	"strings"

	"github.com/canonical/rt-conf/src/data"
	cmp "github.com/canonical/rt-conf/src/ui/components"
	"github.com/canonical/rt-conf/src/ui/config"
	"github.com/canonical/rt-conf/src/ui/styles"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	// keys *listKeyMap
	// TODO: reporpuse this for adding IRQ affinity entries
	// itemGenerator *menuItems
	// delegateKeys *selectKeyMap
	Width  int
	Height int
	iConf  data.InternalConfig
	// for the kernel command line view
	// For the IRQ tunning view
	// irqInputs     []textinput.Model
	// irqFocusIndex int
	Nav        *cmp.MenuNav
	cursorMode cursor.Mode
	errorMsg   string
	logMsg     []string
	renderLog  bool

	// The keymap is consistent across all menus
	currMenu tea.Model

	main MainMenuModel
	irq  IRQMenuModel
	kcmd KcmdlineMenuModel
}

func (m *Model) GetActiveMenu() tea.Model {
	//** Note it's very important `m` be a pointer receiver

	menu := map[config.Views]tea.Model{
		config.INIT_VIEW_ID:            &m.main,
		config.KCMD_VIEW_ID:            &m.kcmd,
		config.KCMD_CONCLUSION_VIEW_ID: &m.kcmd.concl,
		config.IRQ_VIEW_ID:             &m.irq,
		config.IRQ_ADD_EDIT_VIEW_ID:    &m.irq.irq,
		config.IRQ_CONCLUSION_VIEW_ID:  &m.irq.concl,
	}
	mm, ok := menu[m.Nav.GetCurrMenu()]
	if !ok {
		log.Println("ERROR: Menu not found, index: ", m.Nav.GetCurrMenu())
		return &m.main
	}
	return mm
}

func NewModel(c *data.InternalConfig) Model {
	mainMenu := NewMainMenuModel()
	irqMenu := newModelIRQMenuModel()
	kcmd := newKcmdMenuModel(c)

	nav := cmp.GetMenuNavInstance()
	return Model{
		Nav: nav,

		iConf: *c,
		main:  mainMenu,
		irq:   irqMenu,
		kcmd:  kcmd,

		// keys:     listKeys,
		errorMsg: strings.Repeat("\n", len(validationErrorsKcmd)),

		// itemGenerator: &menuOpts,
		cursorMode: cursor.CursorBlink,
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// log.Println("(MAIN UPDATE) Current menu: ", config.Menu[m.Nav.GetCurrMenu()])
	// log.Println("MENU STACK: ", m.Nav.PrintMenuStack())
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		// ** NOTE: This needs to be done here in the main Update() function
		// ** because for some reason if passed for the other
		// ** Update() functions the size of the window is not updated
		// TODO: check why this is happening
		h, v := styles.AppStyle.GetFrameSize()
		m.Width = msg.Width - h
		m.Height = msg.Height - v
		// m.kcmd.list.SetSize(msg.Width-h, msg.Height-v) // TODO: implement this

		// Main Menu
		m.main.list.SetSize(msg.Width-h, msg.Height-v)

		// Kcmdline
		m.kcmd.help.Width = msg.Width
		m.kcmd.Width = msg.Width - h
		m.kcmd.Height = msg.Height - v
		m.kcmd.concl.Height = msg.Height - v
		m.kcmd.concl.Width = msg.Width - h

		// IRQ
		m.irq.list.SetSize(msg.Width-h, msg.Height-v)
		m.irq.Width = msg.Width - h
		m.irq.Height = msg.Height - v
		m.irq.irq.width = msg.Width - h
		m.irq.irq.height = msg.Height - v
		m.irq.concl.Width = msg.Width - h
		m.irq.concl.Height = msg.Height - v

	case tea.KeyMsg:
		// log.Printf("Key pressed: %v", msg.String())
		switch {
		case key.Matches(msg, m.main.keys.Back):
			m.Nav.PrevMenu()
			return m, nil
		// This is genertic for all menus

		case key.Matches(msg, m.main.keys.Quit):
			if msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
			return m, nil

		case key.Matches(msg, m.main.keys.Home):
			m.Nav.ReturnToMainMenu()

		}
	}

	activeMenu := m.GetActiveMenu()
	var cmd tea.Cmd
	// Calling the Update() function of the active menu
	m.currMenu, cmd = activeMenu.Update(msg)
	cmds = append(cmds, cmd)

	// TODO: improve this workaround
	// This is an excecption for the IRQ Add/Editview so the 'q' key can be used
	if m.Nav.GetCurrMenu() == config.IRQ_ADD_EDIT_VIEW_ID {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q":
				return m, tea.Batch(cmds...)
			}
		}
	}

	// Update main menu list view
	newListModel, cmd := m.main.list.Update(msg)
	m.main.list = newListModel
	cmds = append(cmds, cmd)

	// Update Kcmdline list view // TODO: implement this
	// newKcmd, cmd := m.kcmd.list.Update(msg)
	// m.kcmd.list = newKcmd
	// cmds = append(cmds, cmd)

	// Update IRQ list view
	new, cmd := m.irq.list.Update(msg)
	m.irq.list = new
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	// log.Println("\n------------- VIEW -------------------")
	// log.Printf("(VIEW) Current menu: %s", config.Menu[m.Nav.GetCurrMenu()])
	return m.GetActiveMenu().View()
}
