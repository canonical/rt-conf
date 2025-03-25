package ui

import (
	"log"
	"strings"

	"github.com/canonical/rt-conf/src/model"
	cmp "github.com/canonical/rt-conf/src/ui/components"
	"github.com/canonical/rt-conf/src/ui/config"
	"github.com/canonical/rt-conf/src/ui/styles"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

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

type Model struct {
	// keys *listKeyMap
	// TODO: reporpuse this for adding IRQ affinity entries
	// itemGenerator *menuItems
	// delegateKeys *selectKeyMap
	Width  int
	Height int
	iConf  model.InternalConfig
	// for the kernel command line view
	// For the IRQ tuning view
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

type MainMenuModel struct {
	Nav          *cmp.MenuNav
	Width        int
	Height       int
	keys         *KeyMap
	list         list.Model
	delegateKeys *selectKeyMap
}

type IRQMenuModel struct {
	Nav      *cmp.MenuNav
	Width    int
	Height   int
	Index    int
	newEntry bool
	editMode bool
	rules    []model.IRQTuning
	keys     *irqKeyMap
	list     list.Model
	// help     help.Model

	concl IRQConclusion
	irq   IRQAddEditMenu
}

type KcmdlineMenuModel struct {
	Nav        *cmp.MenuNav // Menu Navigation instance
	keys       *kcmdKeyMap
	help       help.Model
	Inputs     []textinput.Model
	concl      KcmdlineConclusion
	Width      int
	Height     int
	FocusIndex int
	errorMsg   string
	iConf      model.InternalConfig
	// keys     *listKeyMap
}

type KcmdlineConclusion struct {
	Nav       *cmp.MenuNav // Menu Navigation instance
	keys      *kcmdKeyMap
	Width     int
	Height    int
	logMsg    []string
	renderLog bool
}

type IRQAddEditMenu struct {
	Nav        *cmp.MenuNav // Menu Navigation instance
	width      int
	height     int
	FocusIndex int
	Inputs     []textinput.Model
	help       help.Model
	keys       *irqAddEditKeyMap
	// errVal     []ErrValidation //TODO: implement this
	errorMsg string

	editMode  bool // false for new entry, true for edit existing entry
	editIndex int  // index of the rule to edit
}

type IRQConclusion struct {
	// The number of IRQs that are being applied to the system
	num    int
	Nav    *cmp.MenuNav // Menu Navigation instance
	keys   *irqKeyMap
	Width  int
	Height int
	logMsg string
	errMsg string
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
	m[0].Focus()
	t.Prompt = config.PrefixCpuList // "CPU Range > "
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
	menuList.KeyMap.NextPage.SetEnabled(false)
	menuList.KeyMap.PrevPage.SetEnabled(false)
	menuList.KeyMap.Filter.SetEnabled(false)

	nav := cmp.GetMenuNavInstance()
	return MainMenuModel{
		Nav:          nav,
		keys:         keys,
		list:         menuList,
		delegateKeys: delegateKeys,
	}
}

func NewModel(c *model.InternalConfig) Model {
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

func newModelIRQMenuModel() IRQMenuModel {
	keys := irqMenuListKeyMap()
	listKeys := irqMenuListKeyMap()

	irq := newIRQAddEditMenuModel()
	concl := newIRQConclusionModel()
	items := []list.Item{
		IRQAffinityRule{filter: config.PrefixIRQFilter,
			cpulist: config.PrefixCpuList},
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedDesc = styles.SelectedListItem
	delegate.Styles.SelectedTitle = styles.SelectedListItem

	delegate.Styles.NormalDesc = styles.ListItem
	delegate.Styles.NormalTitle = styles.ListItem

	delegate.Styles.DimmedDesc = styles.DimmedListItem
	delegate.Styles.DimmedTitle = styles.DimmedListItem

	delegate.Styles.FilterMatch = styles.FilterMatch

	m := list.New(items, delegate, 0, 0)
	m.Title = "IRQ Affinity"
	m.Styles.Title = styles.TitleStyle

	m.KeyMap.Filter.SetEnabled(false)
	m.KeyMap.CursorUp.SetEnabled(false)
	m.KeyMap.CursorDown.SetEnabled(false)
	m.KeyMap.NextPage.SetEnabled(false)
	m.KeyMap.PrevPage.SetEnabled(false)
	m.KeyMap.GoToStart.SetEnabled(false)
	m.KeyMap.GoToEnd.SetEnabled(false)

	m.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.Up,
			listKeys.Down,
			listKeys.Apply,
			listKeys.Add,
			listKeys.Remove,
			listKeys.Back,
		}
	}

	m.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.Up,
			listKeys.Down,
			listKeys.Apply,
			listKeys.Add,
			listKeys.Remove,
			listKeys.Back,
			listKeys.goHome,
		}
	}

	nav := cmp.GetMenuNavInstance()
	return IRQMenuModel{
		Nav:   nav,
		list:  m,
		keys:  keys,
		irq:   irq,
		concl: concl,
		rules: []model.IRQTuning{
			{ // Insert an empty default IRQ tuning rule
				CPUs: "0-n",
				Filter: model.IRQFilter{
					Actions:  "",
					ChipName: "",
					Name:     "",
					Type:     "",
				},
			},
		},
	}
}

func newKcmdMenuModel(c *model.InternalConfig) KcmdlineMenuModel {
	help := help.New()
	inputs := newKcmdTextInputs()
	keys := newkcmdMenuListKeyMap()
	nav := cmp.GetMenuNavInstance()
	concl := newKcmdConclusionModel()
	return KcmdlineMenuModel{
		iConf:  *c,
		Nav:    nav,
		keys:   keys,
		help:   help,
		Inputs: inputs,
		concl:  concl,
	}
}

func newIRQAddEditMenuModel() IRQAddEditMenu {
	help := help.New()
	inputs := newIRQtextInputs()
	keys := irqAddEditListKeyMap()
	nav := cmp.GetMenuNavInstance()
	return IRQAddEditMenu{
		Nav:    nav,
		keys:   keys,
		help:   help,
		Inputs: inputs,
		// Initialize errors strings with empty new line
		// beceusae these will be part of a vertical composed view
		errorMsg: "\n\n",
	}
}

func newKcmdConclusionModel() KcmdlineConclusion {
	keys := newkcmdMenuListKeyMap()
	nav := cmp.GetMenuNavInstance()
	return KcmdlineConclusion{
		Nav:  nav,
		keys: keys,
	}
}

func newIRQConclusionModel() IRQConclusion {
	keys := irqMenuListKeyMap()
	nav := cmp.GetMenuNavInstance()
	return IRQConclusion{
		Nav:  nav,
		keys: keys,
	}
}
