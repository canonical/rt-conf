package ui

import (
	"log"
	"strings"

	"github.com/canonical/rt-conf/src/data"
	"github.com/canonical/rt-conf/src/interrupts"
	cmp "github.com/canonical/rt-conf/src/ui/components"
	"github.com/canonical/rt-conf/src/ui/config"
	"github.com/canonical/rt-conf/src/ui/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type IRQMenuModel struct {
	Nav      *cmp.MenuNav
	Width    int
	Height   int
	Index    int
	newEntry bool
	editMode bool
	rules    []data.IRQTunning
	keys     *irqKeyMap
	list     list.Model
	// help     help.Model

	concl IRQConclusion
	irq   IRQAddEditMenu
}

type IRQRuleMsg struct {
	IRQAffinityRule
}

type IRQAffinityRule struct {
	edited          bool // false means it's a new rule, true means new rule
	index           int
	rule            data.IRQTunning
	filter, cpulist string
}
type irqKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	goHome key.Binding
	Select key.Binding
	Help   key.Binding
	Quit   key.Binding
	Back   key.Binding
	Add    key.Binding
	Remove key.Binding
	Apply  key.Binding
	// Left   key.Binding
	// Right  key.Binding
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k irqKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Apply, k.Add, k.Remove, k.Back, k.Select, k.Help},
		{k.Up, k.Down, k.goHome, k.Select, k.Quit},
	}
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k irqKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Apply, k.Add, k.Remove, k.Back, k.goHome, k.Help}
}

func irqMenuListKeyMap() *irqKeyMap {
	return &irqKeyMap{
		Add: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add item"),
		),
		Remove: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp("x", "remove item"),
		),
		Apply: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "Apply changes on system"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		// Left: key.NewBinding(
		// 	key.WithKeys("left"),
		// 	key.WithHelp("←", "move left"),
		// ),
		// Right: key.NewBinding(
		// 	key.WithKeys("right"),
		// 	key.WithHelp("→", "move right"),
		// ),
		goHome: key.NewBinding(
			key.WithKeys("home"),
			key.WithHelp("home", "main menu"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter", "space"),
			key.WithHelp("enter", "open item"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
	}
}

type StartNewIRQAffinityRule struct {
	editMode bool
	// ** If editMode is true, this is the index of the rule
	// **  to edit as well the cpu list and filter
	index           int
	filter, cpulist string
}

func (i IRQAffinityRule) Title() string       { return i.filter }
func (i IRQAffinityRule) Description() string { return i.cpulist }

// This needs to be implemented to satisfy the list.Item interface
func (i IRQAffinityRule) FilterValue() string { return i.filter }

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
		rules: []data.IRQTunning{
			{ // Insert an empty default IRQ tunning rule
				CPUs: "0-n",
				Filter: data.IRQFilter{
					Actions:  "",
					ChipName: "",
					Name:     "",
					Type:     "",
				},
			},
		},
	}
}

func (m IRQMenuModel) Init() tea.Cmd { return nil }

const (
	irqFilterIndex = iota
	cpuListIndex
	addBtnIndex
	cancelBtnIndex
)

var placeholders_irq = []string{
	config.IrqFilterPlaceholder,
	config.CpuListPlaceholder,
}

func (m *IRQMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case IRQRuleMsg:
		log.Println("<><><><>(RECEIVED) IRQ Affinity Rule: <><><><>")
		irqRule := msg.IRQAffinityRule
		// items := m.list.Items()
		newFilter := config.PrefixIRQFilter + irqRule.filter
		newCpuList := config.PrefixCpuList + irqRule.cpulist
		// log.Println("->>NewFilter:  ", newFilter)
		// log.Println("->>NewCpuList:  ", newCpuList)

		var cmd tea.Cmd
		item := IRQAffinityRule{
			rule:    irqRule.rule,
			filter:  newFilter,
			cpulist: newCpuList,
		}
		if msg.edited {
			// log.Println("->> Edit mode")
			// log.Println("->> Index: ", irqRule.index)
			cmd = m.list.SetItem(irqRule.index, item)
			m.rules[irqRule.index] = irqRule.rule
		} else {
			index := len(m.list.Items())
			cmd = m.list.InsertItem(index, item)
			m.rules = append(m.rules, irqRule.rule)
		}
		return m, cmd

	case tea.KeyMsg:
		// log.Printf("key pressed: %v", msg.String())

		switch {

		// case key.Matches(msg, m.keys.Help):
		// 	log.Println("->>> Help key pressed <<<-")

		case key.Matches(msg, m.keys.Apply):
			// log.Println("-> Entry: Apply changes")
			num := len(m.rules)
			if len(m.rules) == len(m.list.Items()) {
				log.Println(">>> The size is coeherent")
			} else {
				log.Println("[ERROR] >>> The size is NOT coeherent")
				// panic("The size is NOT coeherent")
				// os.Exit(1)
			}
			if num > 0 {
				m.concl.num = num
				m.concl.logMsg =
					"Total IRQ affinity rules applied to the system: "

				var cfg data.InternalConfig
				cfg.Data.Interrupts = m.rules

				err := interrupts.ApplyIRQConfig(&cfg)
				if err != nil {
					m.concl.num = 0
					if strings.Contains(
						err.Error(), "no IRQs matched the filter") {
						m.concl.errMsg =
							"No IRQs matched the given filter(s)."
					} else {
						m.concl.errMsg =
							"ERROR: " + err.Error()
					}
					log.Println("ERROR: ", err)
				}
				m.Nav.SetNewMenu(config.IRQ_CONCLUSION_VIEW_ID)
			}

		case key.Matches(msg, m.keys.Add):
			// log.Println("-> Entry: Add new IRQ rule")
			// TODO: implement cleanup of AddEditMenu fields
			cmd := func() tea.Msg {
				return StartNewIRQAffinityRule{editMode: false}
			}
			m.Nav.SetNewMenu(config.IRQ_ADD_EDIT_VIEW_ID)
			return m, cmd

		case key.Matches(msg, m.keys.Remove):
			// log.Println("-> Entry: Remove IRQ rule")
			idx := m.list.Index()
			m.rules = append(m.rules[:idx], m.rules[idx+1:]...)
			m.list.RemoveItem(idx)

		case key.Matches(msg, m.irq.keys.Select):
			// log.Printf("confirmed select pressed")
			// log.Printf("previous menu: %v", config.Menu[m.Nav.GetCurrMenu()])
			cpulist := strings.TrimPrefix(
				m.list.SelectedItem().(IRQAffinityRule).cpulist,
				config.PrefixCpuList,
			)
			filter := strings.TrimPrefix(
				m.list.SelectedItem().(IRQAffinityRule).filter,
				config.PrefixIRQFilter,
			)
			cmd := func() tea.Msg {
				return StartNewIRQAffinityRule{
					editMode: true,
					index:    m.list.Index(),
					cpulist:  cpulist,
					filter:   filter,
				}
			}
			m.Nav.SetNewMenu(config.IRQ_ADD_EDIT_VIEW_ID)

			// log.Println("Current menu: ", config.Menu[m.Nav.GetCurrMenu()])
			return m, cmd

		}
	}
	return m, nil
}

func (m IRQMenuModel) View() string {
	m.list.KeyMap.Filter.SetEnabled(false)
	m.list.KeyMap.CursorUp.SetEnabled(false)
	m.list.KeyMap.CursorDown.SetEnabled(false)
	m.list.KeyMap.NextPage.SetEnabled(false)
	m.list.KeyMap.PrevPage.SetEnabled(false)
	m.list.KeyMap.GoToStart.SetEnabled(false)
	m.list.KeyMap.GoToEnd.SetEnabled(false)
	m.list.KeyMap.Filter.SetEnabled(false)
	return styles.AppStyle.Render(m.list.View())
}
