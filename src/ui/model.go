package ui

import (
	"log"

	"github.com/canonical/rt-conf/src/data"
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
	goHome     key.Binding
	selectMenu key.Binding
	// toggleSpinner    key.Binding
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
		// toggleSpinner: key.NewBinding(
		// 	key.WithKeys("s"),
		// 	key.WithHelp("s", "toggle spinner"),
		// ),
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
}

func NewModel(c *data.InternalConfig) Model {
	var (
		itemGenerator menuItems
		delegateKeys  = newDelegateKeyMap()
		listKeys      = newListKeyMap()
	)

	// Make initial list of items
	const numItems = 2
	items := make([]list.Item, numItems)
	for i := 0; i < numItems; i++ {
		items[i] = itemGenerator.next()
	}

	// Setup list
	delegate := newItemDelegate(delegateKeys)
	menuList := list.New(items, delegate, 0, 0)
	menuList.Title = "rt-conf tool"
	menuList.Styles.Title = titleStyle
	menuList.SetShowStatusBar(false)

	menuList.AdditionalFullHelpKeys = func() []key.Binding {

		return []key.Binding{
			// listKeys.toggleSpinner,
			// listKeys.insertItem,
			// listKeys.toggleTitleBar,
			// listKeys.togglePagination,
			listKeys.goHome,
			listKeys.toggleHelpMenu,
		}
	}

	return Model{
		inputs:        newTextInputs(),
		iconf:         *c,
		list:          menuList,
		keys:          listKeys,
		delegateKeys:  delegateKeys,
		itemGenerator: &itemGenerator,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

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
		if m.list.FilterState() == list.Filtering {
			break
		}

		if m.currMenu == mainMenu {
			switch {
			// case key.Matches(msg, m.keys.toggleSpinner):
			// 	cmd := m.list.ToggleSpinner()
			// 	return m, cmd

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
				return m, nil
			// case key.Matches(msg, m.keys.toggleTitleBar):
			// 	v := !m.list.ShowTitle()
			// 	m.list.SetShowTitle(v)
			// 	m.list.SetShowFilter(v)
			// 	m.list.SetFilteringEnabled(v)
			// 	return m, nil

			// case key.Matches(msg, m.keys.togglePagination):
			// 	m.list.SetShowPagination(!m.list.ShowPagination())
			// 	return m, nil

			case key.Matches(msg, m.keys.toggleHelpMenu):
				m.list.SetShowHelp(!m.list.ShowHelp())
				return m, nil

				// case key.Matches(msg, m.keys.insertItem):
				// 	m.delegateKeys.choose.SetEnabled(true)
				// 	newItem := m.itemGenerator.next()
				// 	insCmd := m.list.InsertItem(0, newItem)
				// 	statusCmd := m.list.NewStatusMessage(statusMessageStyle("Added " + newItem.Title()))
				// 	return m, tea.Batch(insCmd, statusCmd)
			}
		} else {

			// if m.currMenu == kcmdlineMenu && m.prevMenu == mainMenu {
			// 	// Cleanup the inputs
			// 	log.Println("Cleaning up inputs")
			// 	m.inputs = nil
			// }

			switch msg.String() {
			case "ctrl+c", "esc":
				return m, tea.Quit

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

				// Did the user press enter while the submit button was focused?
				// If so, exit.
				if s == "enter" && m.focusIndex == len(m.inputs) {
					return m, tea.Quit
				}

				log.Println("Focus index is: ", m.focusIndex)
				log.Println("Value: ", m.inputs[m.focusIndex].Value())

				m.Validation() //TODO: VALIDATION HAPPENS HERE

				// Cycle indexes
				if s == "up" || s == "shift+tab" {
					m.focusIndex--
				} else {
					m.focusIndex++
				}

				if m.focusIndex > len(m.inputs) {
					m.focusIndex = 0
				} else if m.focusIndex < 0 {
					m.focusIndex = len(m.inputs)
				}

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
			cmdInput = m.updateInputs(msg)
		}
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)
	cmds = append(cmds, cmdInput)

	return m, tea.Batch(cmds...)
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)

		// log.Println("Value is: ", m.inputs[i].Value())
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
