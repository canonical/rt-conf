package ui

import (
	"strings"

	"github.com/canonical/rt-conf/src/data"
	"github.com/canonical/rt-conf/src/ui/styles"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
)

// TODO: Fix menu navigation
// TODO: Fix inner menu help view

type item struct {
	title       string
	description string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

func InitNewIRQTextInputs() []textinput.Model {
	m := newIRQtextInputs()
	m[0].Focus()
	m[0].PromptStyle = styles.FocusedStyle
	m[0].TextStyle = styles.FocusedStyle
	m[0].Placeholder = irqFilterPlaceholder

	return m
}

func newIRQtextInputs() []textinput.Model {
	m := make([]textinput.Model, 2)
	t := textinput.New()
		t.Cursor.Style = styles.CursorStyle
	t.Cursor.SetMode(cursor.CursorBlink) // TODO: check why this isn't working
			t.CharLimit = 64

	// TODO: This order needs to be reviwed
	t.Prompt = "Filter > "
	m[0] = t
	t.Prompt = "CPU Range > "
	m[1] = t

	return m
}

func newKcmdTextInputs() []textinput.Model {
	m := make([]textinput.Model, 5)

	var t textinput.Model
	for i := range m {
		t = textinput.New()
		t.Cursor.Style = styles.CursorStyle
		t.CharLimit = 32
		// TODO: check why Cursor isn't blinking
		t.Cursor.SetMode(cursor.CursorBlink)

		switch i {
		case isolcpusIndex:
			t.Prompt = "Isolate CPUs from general execution (isolcpus) > "
			/* The placeholder is necessary only in the first, because the
			dynamic placeholders start to work after the first
			move of the cursor (either to up or down) */
			// TODO: investigate the dynamic placeholder refresh
			t.Placeholder = cpuListPlaceholder
			t.Focus()
			t.PromptStyle = styles.FocusedStyle
			t.TextStyle = styles.FocusedStyle
		case nohzIndex:
			t.Prompt = "Enable dyntick mode (nohz) > "
			t.CharLimit = 3
		case nohzFullIndex:
			t.Prompt = "Adaptive ticks CPUs (nohz_full) > "
		case kthreadsCPUsIndex:
			t.Prompt = "CPUs to handle kernel threads (kthread_cpus) > "
		case irqaffinityIndex:
			t.Prompt = "CPUs to handle IRQs (irqaffinity) > "
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
	iConf         data.InternalConfig
	// for the kernel command line view
	kcmdInputs     []textinput.Model
	kcmdFocusIndex int
	// For the IRQ tunning view
	irqInputs     []textinput.Model
	irqFocusIndex int
	cursorMode    cursor.Mode
	prevMenu      menuOpt
	currMenu      menuOpt
	errorMsg      string
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

	return Model{
		// TODO: Fix this info msg, put in a better place
		// logMsg:        logmsg[:],
		kcmdInputs:    newKcmdTextInputs(),
		irqInputs:     InitNewIRQTextInputs(),
		help:          help.New(), // TODO: Check NEED for custom style
		iConf:         *c,
		list:          menuList,
		keys:          listKeys,
		errorMsg:      strings.Repeat("\n", len(validationErrors)),
		delegateKeys:  delegateKeys,
		itemGenerator: &menuOpts,
		cursorMode:    cursor.CursorBlink,
	}
}
