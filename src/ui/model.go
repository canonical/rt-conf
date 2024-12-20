package ui

import (
	"log"
	"strconv"

	"github.com/canonical/rt-conf/src/data"
	"github.com/canonical/rt-conf/src/kcmd"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// TODO: Refact this HORRIBLE code
// TODO: Fix menu navigation
// TODO: Fix inner menu help view

type item struct {
	title       string
	description string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type listKeyMap struct {
	// toggleTitleBar   key.Binding
	// insertItem       key.Binding
	goHome        key.Binding
	selectMenu    key.Binding
	toggleSpinner key.Binding
	// togglePagination key.Binding
	toggleHelpMenu key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		// insertItem: key.NewBinding(
		// 	key.WithKeys("a"),
		// 	key.WithHelp("a", "add item"),
		// ),
		// toggleTitleBar: key.NewBinding(
		// 	key.WithKeys("T"),
		// 	key.WithHelp("T", "toggle title"),
		// ),
		goHome: key.NewBinding(
			key.WithKeys("g", "home"),
			key.WithHelp("esc", "go home"),
		),
		selectMenu: key.NewBinding(
			key.WithKeys("enter", "space"),
			key.WithHelp("enter", "select menu"),
		),
		toggleSpinner: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "toggle spinner"),
		),
		// togglePagination: key.NewBinding(
		// 	key.WithKeys("P"),
		// 	key.WithHelp("P", "toggle pagination"),
		// ),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}

func newTextInputs() []textinput.Model {
	m := make([]textinput.Model, 3)

	var t textinput.Model
	for i := range m {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Prompt = "Isolate CPUs from serving IRQs > "
			t.Placeholder = "2-n"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Prompt = "Enable dyntick mode? (true/false) > "
			t.Placeholder = "true"
			t.CharLimit = 5
		case 2:
			t.Prompt = "Adaptive ticks CPUs > "
			t.Placeholder = "2-n"
			// t.EchoMode = textinput.EchoPassword
			// t.EchoCharacter = '•'
		}

		m[i] = t
	}

	return m
}

type Model struct {
	list          list.Model
	itemGenerator *menuItems
	keys          *listKeyMap
	help          help.Model
	delegateKeys  *selectKeyMap
	width         int
	height        int
	iconf         data.InternalConfig
	inputs        []textinput.Model
	focusIndex    int
	cursorMode    cursor.Mode
	prevMenu      menuOpt
	currMenu      menuOpt
	errorMsg      string
	infoMsg       string
}

func NewModel(c *data.InternalConfig) Model {
	var (
		menuOpts     menuItems
		delegateKeys = newDelegateKeyMap()
		listKeys     = newListKeyMap()
	)

	menuOpts.Init()
	// Make initial list of items
	items := make([]list.Item, menuOpts.Size())
	for i := 0; i < menuOpts.Size(); i++ {
		items[i] = menuOpts.next()
	}

	// Setup list
	delegate := newItemDelegate(delegateKeys)
	menuList := list.New(items, delegate, 0, 0)
	menuList.Title = "rt-conf tool"
	menuList.Styles.Title = titleStyle
	menuList.SetShowStatusBar(false)

	menuList.AdditionalFullHelpKeys = func() []key.Binding {

		return []key.Binding{
			listKeys.toggleSpinner,
			// listKeys.insertItem,
			// listKeys.toggleTitleBar,
			// listKeys.togglePagination,
			listKeys.goHome,
			listKeys.toggleHelpMenu,
		}
	}

	return Model{
		// TODO: Fix this info msg, put in a better place
		infoMsg:       "Please fill all fields before submit\n",
		inputs:        newTextInputs(),
		iconf:         *c,
		list:          menuList,
		keys:          listKeys,
		delegateKeys:  delegateKeys,
		itemGenerator: &menuOpts,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// TODO: FOCUS ON THE [ APPLY ] FUNCTIONALITY
// TODO:

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmdInput tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.width = msg.Width - h
		m.height = msg.Height - v

		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.help.Width = msg.Width

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.

		if m.currMenu != mainMenu {

			// TODO: Move this logic to a func specific for kcmdlineMenu
			switch msg.String() {
			case "ctrl+c", "esc":
				return m, tea.Quit

			case "g":
				m.currMenu = mainMenu
				return m, nil

			// Change cursor mode
			case "ctrl+r":
				m.cursorMode++
				if m.cursorMode > cursor.CursorHide {
					m.cursorMode = cursor.CursorBlink
				}
				cmds := make([]tea.Cmd, len(m.inputs))
				for i := range m.inputs {
					cmds[i] = m.inputs[i].Cursor.SetMode(m.cursorMode)
				}
				return m, tea.Batch(cmds...)

			// Set focus to next input
			case "tab", "shift+tab", "enter", "up", "down":
				s := msg.String()

				log.Println("focusIndex on Update: ", m.focusIndex)

				isValid := m.Validation()
				log.Println("isValid: ", isValid)
				// Did the user press enter while the submit button was focused?
				// If so, exit.
				if s == "enter" && m.focusIndex == len(m.inputs) {
					log.Println("Apply changes")

					valid := m.Validation()

					dyntickIdle, err := strconv.ParseBool(m.inputs[enableDynticks].Value())
					if err != nil {
						log.Printf("Failed to parse dyntick idle value: %v", err)
						break
					}
					if !valid {
						break
					}
					// TODO: fix this
					m.iconf.Data.KernelCmdline.DyntickIdle = dyntickIdle
					m.iconf.Data.KernelCmdline.IsolCPUs = m.inputs[isolatecpus].Value()
					m.iconf.Data.KernelCmdline.AdaptiveTicks = m.inputs[enableDynticks].Value()

					// IsolCPUs      string `yaml:"isolcpus"`       // Isolate CPUs
					// DyntickIdle   bool   `yaml:"dyntick-idle"`   // Enable/Disable dyntick idle
					// AdaptiveTicks string `yaml:"adaptive-ticks"` // CPUs for adaptive ticks

					kcmd.ProcessKcmdArgs(&m.iconf)
					// TODO: Apply the changes call the kcmdline funcs
					// TODO: Present a new view (maybe a new block with text)
				}

				log.Println("Focus index is: ", m.focusIndex)
				if m.focusIndex < len(m.inputs) {
					log.Println("Value: ", m.inputs[m.focusIndex].Value())
				}

				// Cycle indexes
				if s == "up" || s == "shift+tab" {
					m.PrevIndex()
				} else {
					m.NextIndex()
				}

				// if m.focusIndex > len(m.inputs) {
				// 	m.focusIndex = 0
				// } else if m.focusIndex < 0 {
				// 	m.focusIndex = len(m.inputs)
				// }

				cmds := make([]tea.Cmd, len(m.inputs))
				for i := 0; i <= len(m.inputs)-1; i++ {
					if i == m.focusIndex {
						// Set focused state
						cmds[i] = m.inputs[i].Focus()
						m.inputs[i].PromptStyle = focusedStyle
						m.inputs[i].TextStyle = focusedStyle
						continue
					}
					// Remove focused state
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = noStyle
					m.inputs[i].TextStyle = noStyle
				}

				return m, tea.Batch(cmds...)
			}
		}
	}

	cmdMainMenu := m.updateMainMenu(msg)
	cmdInput = m.updateInputs(msg)
	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)
	cmds = append(cmds, cmdInput)
	cmds = append(cmds, cmdMainMenu)

	return m, tea.Batch(cmds...)
}

func (m *Model) NextIndex() {
	m.focusIndex++
	if m.focusIndex > len(m.inputs) {
		m.focusIndex = 0
	}
}

func (m *Model) PrevIndex() {
	m.focusIndex--
	if m.focusIndex == -1 {
		m.focusIndex = len(m.inputs)
	}
}

// func (m *Model) updateKcmdlineMenu(msg tea.Msg) tea.Cmd {
// 	cmds := make([]tea.Cmd, len(m.inputs))
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch {
// 		case key.Matches(msg, m.keys.goHome):
// 			m.currMenu = mainMenu
// 		}
// 	return tea.Batch(cmds...)
// }

func (m *Model) updateMainMenu(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	if m.currMenu == mainMenu {
		switch msg := msg.(type) {
		case tea.KeyMsg:

			if m.list.FilterState() == list.Filtering {
				break
			}
			switch {
			case key.Matches(msg, m.keys.toggleSpinner):
				cmd := m.list.ToggleSpinner()
				cmds = append(cmds, cmd)

			case key.Matches(msg, m.keys.goHome):
				m.currMenu = mainMenu

			case key.Matches(msg, m.keys.selectMenu):
				selected := m.list.SelectedItem().(item)

				// TODO: Improve this selection logic
				// It could be indexed by the menu item
				switch selected.Title() {
				case "Kernel cmdline":
					m.prevMenu = mainMenu
					m.currMenu = kcmdlineMenu
				case "IRQ Affinity":
					m.prevMenu = mainMenu
					m.currMenu = irqAffinityMenu
				}

			case key.Matches(msg, m.keys.toggleHelpMenu):
				m.list.SetShowHelp(!m.list.ShowHelp())
				// return m, nil
			}

		}

	}

	return tea.Batch(cmds...)
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	if m.currMenu == kcmdlineMenu {
		// Only text inputs with Focus() set will respond, so it's safe to simply
		// update all of them here without any further logic.
		for i := range m.inputs {
			m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
		}
	}

	return tea.Batch(cmds...)
}

// TODO: Need to think a way to model the navigation between menus
func (m Model) View() string {
	switch m.currMenu {
	case kcmdlineMenu:
		return appStyle.Render(m.kcmdlineView())
	case irqAffinityMenu:
		return appStyle.Render(m.irqAffinityView())
	default:
		return appStyle.Render(m.list.View())
	}
}
