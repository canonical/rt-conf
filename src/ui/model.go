package ui

import (
	"strings"

	"github.com/canonical/rt-conf/src/data"
	cmp "github.com/canonical/rt-conf/src/ui/components"
	"github.com/canonical/rt-conf/src/ui/config"
	"github.com/canonical/rt-conf/src/ui/styles"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
)

type Model struct {
	// keys *listKeyMap
	// TODO: reporpuse this for adding IRQ affinity entries
	// itemGenerator *menuItems
	// delegateKeys *selectKeyMap
	width  int
	height int
	iConf  data.InternalConfig
	// for the kernel command line view
	// For the IRQ tunning view
	// irqInputs     []textinput.Model
	// irqFocusIndex int
	nav        *cmp.MenuNav
	cursorMode cursor.Mode
	errorMsg   string
	logMsg     []string
	renderLog  bool

	// The keymap is consistent across all menus

	main MainMenuModel
	irq  IRQMenuModel
	kcmd KcmdlineMenuModel
}

type MainMenuModel struct {
	keys         *KeyMap
	list         list.Model
	delegateKeys *selectKeyMap
}

type IRQMenuModel struct {
	width    int
	height   int
	Index    int
	newEntry bool
	editMode bool
	keys     *irqKeyMap
	list     list.Model
	help     help.Model

	irq IRQAddEditMenu
	// keys     *listKeyMap
	// nav components.Navigation
}

type KcmdlineMenuModel struct {
	keys       *kcmdKeyMap
	help       help.Model
	Inputs     []textinput.Model
	FocusIndex int
	// nav components.Navigation
}

type IRQAddEditMenu struct {
	FocusIndex     int
	Inputs         []textinput.Model
	help           help.Model
	keys           *irqKeyMap
	errorMsgFilter string
	errorMsgCpu    string
}

// TODO: Fix inner menu help view

type menuItem struct {
	title       string
	description string
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.description }
func (i menuItem) FilterValue() string { return i.title }
func (r *menuItem) Size() int          { return config.NUMBER_OF_MENUS }
func (r *menuItem) Init()              {}

func InitNewIRQTextInputs() []textinput.Model {
	m := newIRQtextInputs()
	m[0].Focus()
	m[0].PromptStyle = styles.FocusedStyle
	m[0].TextStyle = styles.FocusedStyle
	m[0].Placeholder = config.IrqFilterPlaceholder

	return m
}

func newIRQtextInputs() []textinput.Model {
	m := make([]textinput.Model, 2)
	t := textinput.New()
	t.Cursor.Style = styles.CursorStyle
	t.Cursor.SetMode(cursor.CursorBlink) // TODO: check why this isn't working
	t.CharLimit = 64

	// TODO: This order needs to be reviewed
	t.Prompt = config.PrefixIRQFilter // "Filter > "
	m[0] = t
	m[0].Placeholder = config.IrqFilterPlaceholder
	t.Prompt = config.PrefixCpuRange // "CPU Range > "
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
			t.Placeholder = config.CpuListPlaceholder
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

func NewMainMenuModel() MainMenuModel {
	keys := newMainMenuListKeyMap()
	var delegateKeys = newDelegateKeyMapMainMenu()

	items := []list.Item{
		menuItem{config.MENU_KCMDLINE, config.DESC_KCMDLINE},
		menuItem{config.MENU_IRQAFFINITY, config.DESC_IRQAFFINITY},
		menuItem{config.MENU_PWRMGMT, config.DESC_KCMDLINE},
	}

	delegate := newItemDelegateMainMenu(delegateKeys)
	delegate.Styles.SelectedDesc = styles.SelectedDesc
	delegate.Styles.SelectedTitle = styles.SelectedTitle

	delegate.Styles.NormalDesc = styles.NormalDesc
	delegate.Styles.NormalTitle = styles.NormalTitle

	delegate.Styles.DimmedDesc = styles.DimmedDesc
	delegate.Styles.DimmedTitle = styles.DimmedTitle

	delegate.Styles.FilterMatch = styles.FilterMatch
	menuList := list.New(items, delegate, 0, 0)
	menuList.SetShowHelp(true)
	menuList.Title = "rt-conf tool"
	menuList.Styles.Title = styles.TitleStyle
	menuList.SetShowStatusBar(false)
	menuList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			keys.Select,
		}
	}

	return MainMenuModel{
		keys:         keys,
		list:         menuList,
		delegateKeys: delegateKeys,
	}
}

func NewModel(c *data.InternalConfig) Model {
	mainMenu := NewMainMenuModel()
	irqMenu := newModelIRQMenuModel()
	kcmd := newKcmdMenuModel()
	nav := cmp.NewMenuNav()

	return Model{
		// TODO: Fix this info msg, put in a better place
		// logMsg:        logmsg[:],
		// irqInputs:  InitNewIRQTextInputs(),

		nav: nav, // TODO: send this as tea.Msg

		// help:  help.New(), // TODO: Check NEED for custom style
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

func newModelIRQMenuModel() IRQMenuModel {
	keys := irqMenuListKeyMap()
	help := help.New()
	irq := newIRQAddEditMenuModel()
	items := []list.Item{
		irqAffinityRule{filter: "Filter > ", cpulist: "CPU List > "},
	}
	// m := list.New(items, newItemDelegate(newDelegateKeyMap()), 0, 0)
	m := list.New(items, list.NewDefaultDelegate(), 0, 0)
	m.Title = "IRQ Affinity"
	m.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			keys.Add,
			keys.Remove,
			keys.Apply,
		}
	}
	return IRQMenuModel{
		list: m,
		keys: keys,
		help: help,
		irq:  irq,
	}
}

func newKcmdMenuModel() KcmdlineMenuModel {
	help := help.New()
	inputs := newKcmdTextInputs()
	keys := newkcmdMenuListKeyMap()
	return KcmdlineMenuModel{
		keys:   keys,
		help:   help,
		Inputs: inputs,
	}
}

func newIRQAddEditMenuModel() IRQAddEditMenu {
	help := help.New()
	inputs := newIRQtextInputs()
	keys := irqMenuListKeyMap()
	return IRQAddEditMenu{
		keys:           keys,
		help:           help,
		Inputs:         inputs,
		errorMsgFilter: "\n",
		errorMsgCpu:    "\n",
	}

}
