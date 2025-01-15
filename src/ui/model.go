package ui

import (
	"github.com/canonical/rt-conf/src/data"
	"github.com/canonical/rt-conf/src/ui/styles"
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

func newTextInputs() []textinput.Model {
	m := make([]textinput.Model, 3)

	var t textinput.Model
	for i := range m {
		t = textinput.New()
		t.Cursor.Style = styles.CursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Prompt = "Isolate CPUs from general execution (isolcpus) > "
			/* The placeholder is necessary only in the first, because the
			dynamic placeholders start to work after the first
			move of the cursor (either to up or down) */
			// TODO: investigate the dynamic placeholder refresh
			t.Placeholder = firstPlaceholder
			t.Focus()
			t.PromptStyle = styles.FocusedStyle
			t.TextStyle = styles.FocusedStyle
		case 1:
			t.Prompt = "Enable dyntick mode (nohz) > "
			t.CharLimit = 1
		case 2:
			t.Prompt = "Adaptive ticks CPUs (nohz_full) > "
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
	menuList.Styles.Title = styles.TitleStyle
	menuList.SetShowStatusBar(false)

	menuList.Help = help.New()

	menuList.AdditionalFullHelpKeys = func() []key.Binding {

		return []key.Binding{
			listKeys.goHome,
			listKeys.Help,
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
		help:          help.New(), // TODO: Check NEED for custom style
		iconf:         *c,
		list:          menuList,
		keys:          listKeys,
		delegateKeys:  delegateKeys,
		itemGenerator: &menuOpts,
		cursorMode:    cursor.CursorBlink,
	}
}

// TODO: figure out what is wrong with this
// * NOTE: For some reason this update function isn't working properlly *
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
