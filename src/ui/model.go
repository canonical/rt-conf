package ui

import (
	"github.com/canonical/rt-conf/src/data"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
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
			key.WithKeys("g", "home", "esc"),
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
			t.Prompt = "Isolate CPUs from general execution > "
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
	logMsg        []string
	renderLog     bool
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

	var logmsg [8]string
	for i := range 8 {
		logmsg[i] = "\n"
	}
	// NOTE: * 8 * because:
	// There is 8 lines of output when processing the kcmdline functions

	return Model{
		// TODO: Fix this info msg, put in a better place
		infoMsg:       "Please fill all fields before submit\n",
		logMsg:        logmsg[:],
		inputs:        newTextInputs(),
		iconf:         *c,
		list:          menuList,
		keys:          listKeys,
		delegateKeys:  delegateKeys,
		itemGenerator: &menuOpts,
	}
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
