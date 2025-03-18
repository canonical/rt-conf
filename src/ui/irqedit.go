package ui

import (
	"fmt"
	"log"
	"strings"

	"github.com/canonical/rt-conf/src/data"
	cmp "github.com/canonical/rt-conf/src/ui/components"
	"github.com/canonical/rt-conf/src/ui/config"
	"github.com/canonical/rt-conf/src/ui/styles"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

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

// TODO: check if there is a better way to handle the keybindings
// ** There is a lot of repetition here

type irqAddEditKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	goHome key.Binding
	Back   key.Binding
	Select key.Binding
	Help   key.Binding
	Quit   key.Binding
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k irqAddEditKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right, k.Back},
		{k.Select, k.goHome, k.Quit, k.Help},
	}
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k irqAddEditKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.goHome, k.Quit, k.Help}
}

func irqAddEditListKeyMap() *irqAddEditKeyMap {
	return &irqAddEditKeyMap{
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", "left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("→", "right"),
		),
		goHome: key.NewBinding(
			key.WithKeys("home"),
			key.WithHelp("home", "Main menu"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
	}
}

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

func (m IRQAddEditMenu) Init() tea.Cmd { return nil }

func (m *IRQAddEditMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(m.Inputs), len(m.Inputs)+4)

	totalIrqItems := len(m.Inputs) + 2
	index := cmp.NewNavigation(&m.FocusIndex, &totalIrqItems)
	// log.Println("\tFocus index: ", m.FocusIndex)
	// log.Println("\tAddres: ", &m.FocusIndex)
	// log.Printf("---IRQAddEditMenuUpdate: ")
	// log.Printf("Current menu: %v", config.Menu[m.Nav.GetCurrMenu()])

	// TODO: here it should be implemented the logic to add or edit IRQs
	switch msg := msg.(type) {

	case StartNewIRQAffinityRule:
		// log.Println("<><><> StartNewIRQAffinityRule")
		// log.Println(">>> Edit mode: ", msg.editMode)
		// log.Println(">>> Index: ", msg.index)

		m.editMode = msg.editMode
		m.editIndex = msg.index
		m.FocusIndex = 0
		m.Inputs[irqFilterIndex].Focus()
		m.Inputs[irqFilterIndex].PromptStyle = styles.FocusedStyle
		m.Inputs[irqFilterIndex].PromptStyle = styles.FocusedStyle
		m.Inputs[irqFilterIndex].TextStyle = styles.FocusedStyle
		m.Inputs[irqFilterIndex].Placeholder = placeholders_irq[irqFilterIndex]
		if !msg.editMode {
			m.Inputs[irqFilterIndex].SetValue("")
			m.Inputs[cpuListIndex].SetValue("")
		} else {
			m.Inputs[irqFilterIndex].SetValue(msg.filter)
			m.Inputs[cpuListIndex].SetValue(msg.cpulist)
		}

	case tea.KeyMsg:
		var isValid bool
		switch {
		// TODO: find a way to allow use of 'q' key on the text inputs
		case key.Matches(msg, m.keys.Quit):
			if msg.String() == "ctrl+c" {
				return m, tea.Quit
			}

		case key.Matches(msg, m.keys.Help):
			var cmd tea.Cmd
			if m.FocusIndex < addBtnIndex {
				m.Inputs[m.FocusIndex].Blur()
			}
			m.help.ShowAll = !m.help.ShowAll
			if m.FocusIndex < addBtnIndex {
				cmd = m.Inputs[m.FocusIndex].Focus()
			}
			return m, tea.Sequence(cmd)

		case key.Matches(msg, m.keys.Up):
			m.RunInputValidation()
			index.Prev()

		case key.Matches(msg, m.keys.Down):
			m.RunInputValidation()
			index.Next()

		case key.Matches(msg, m.keys.Left):
			if m.FocusIndex == addBtnIndex {
				index.Prev()
			} else if m.FocusIndex == cancelBtnIndex {
				index.Prev()
			}

		case key.Matches(msg, m.keys.Right):
			if m.FocusIndex == addBtnIndex {
				index.Next()
			} else if m.FocusIndex == cancelBtnIndex {
				index.Next()
			}

		case key.Matches(msg, m.keys.Select):
			isValid = m.AreValidInputs()
			// log.Println(">>> Select key pressed")

			// TODO: fix bug where it's necessary to press twice the buttons
			// Handle [ Cancel ] button
			if m.FocusIndex == cancelBtnIndex {
				// log.Println(">>> Cancel changes")
				m.Nav.PrevMenu()
				break
			}

			if m.FocusIndex < len(m.Inputs) /*addBtnIndex*/ {
				index.Next()
				break
			}

			// Did the user press enter while the [ ADD ] button was focused?
			if m.FocusIndex == addBtnIndex {
				log.Println(">>> Apply changes | isValid: ", isValid)

				var empty int
				for i := range m.Inputs {
					v := m.Inputs[i].Value()
					if v == "" {
						empty++
					}
				}
				if empty == len(m.Inputs) {
					isValid = false
					m.errorMsg = "\nAll fields are empty\n"
					break
				}
				if empty > 0 {
					isValid = false
					m.errorMsg = "\nA field is empty\n"
					break
				}

				if !isValid {
					break
				}

				rawFilter := m.Inputs[irqFilterIndex].Value()
				cpulist := m.Inputs[cpuListIndex].Value()

				irqFilter, _ := ParseIRQFilter(rawFilter)
				// if err != nil {
				// 	m.errorMsg = "ERROR: " + err.Error() + "\n"
				// 	break
				// }

				irqAffinityRule := IRQAffinityRule{
					rule: data.IRQTunning{
						Filter: irqFilter,
						CPUs:   cpulist,
					},
					filter:  rawFilter,
					cpulist: cpulist,
					edited:  m.editMode,
					index:   m.editIndex,
				}
				cmd := func() tea.Msg { return IRQRuleMsg{irqAffinityRule} }
				cmds = append(cmds, cmd)

				// TODO once validated, it must:
				// ** 1. Parse the values for the IRQFilter struct
				// ** 1.1 The func for parsing of IRQFilter  must be implemented
				// ** 1.2 The cpuList value can be copy directly

				// ** 2. Create a tea.Cmd message to send the IRQFilter struct
				// ** 2.1 The tea.Cmd must be added into the tea.Batch(cmds...)
				m.Nav.PrevMenu()

			}

		default:
			// log.Printf("key pressed: %v", msg.String())
		}
	}

	for i := 0; i <= len(m.Inputs)-1; i++ {
		if i == m.FocusIndex {
			// Set focused state
			cmds[i] = m.Inputs[i].Focus()
			m.Inputs[i].PromptStyle = styles.FocusedStyle
			m.Inputs[i].TextStyle = styles.FocusedStyle
			m.Inputs[i].Placeholder = placeholders_irq[i]
			continue
		}
		// Remove focused state
		m.Inputs[i].Blur()
		m.Inputs[i].PromptStyle = styles.NoStyle
		m.Inputs[i].TextStyle = styles.NoStyle
		m.Inputs[i].Placeholder = ""
	}

	for i := range m.Inputs {
		m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m IRQAddEditMenu) View() string {
	var s string // the view

	title := styles.InnerMenuStyle("Add IRQ Affinity Rule")
	desc := styles.Section.
		Render("Allocate specific CPUs to IRQs matching the filter.")

	// TODO: bold cpu range and filter
	var b strings.Builder
	for i := range m.Inputs {
		b.WriteString(m.Inputs[i].View())
		if i < len(m.Inputs)-1 {
			b.WriteRune('\n')
		}
	}
	b.WriteString("\n")

	var verticalPadding strings.Builder
	verticalPadding.WriteString("\n\n")

	addBtn := cmp.NewButton("Add")
	cancelBtn := cmp.NewButton("Cancel")

	btns := []*cmp.Button{addBtn, cancelBtn}

	for i, btn := range btns {
		if m.FocusIndex == i+(len(m.Inputs)) {
			btn.SetFocused()
		} else {
			btn.SetBlurred()
		}
	}

	// log.Println("--- (IRQ ADD/EDIT VIEW) m.irq.FocusIndex: ", m.FocusIndex)

	fmt.Fprintf(&b, "\n%s\n",
		styles.JoinHorizontal(
			addBtn.Render(),
			cancelBtn.Render(),
		))

	body := b.String()

	helpView := m.help.View(m.keys)

	height := m.height -
		strings.Count(title, "\n") -
		strings.Count(desc, "\n") -
		strings.Count(helpView, "\n") -
		strings.Count(m.errorMsg, "\n") -
		strings.Count(b.String(), "\n") -
		// verticalPadding is used twice
		strings.Count(verticalPadding.String(), "\n") -
		strings.Count(verticalPadding.String(), "\n")

	// log.Println("--- (IRQ ADD/EDIT VIEW) m.Height: ", m.height)
	// log.Println("--- (IRQ ADD/EDIT VIEW) height: ", height)

	if height < 0 {
		height = 1
	}
	// log.Println("--- recalculated height: ", height)

	s +=
		title +
			verticalPadding.String() +
			desc +
			verticalPadding.String() +
			body +
			strings.Repeat("\n", height/2) +
			styles.ErrorMessageStyle(m.errorMsg) +
			strings.Repeat("\n", height/2) +
			helpView

	return styles.AppStyle.Render(s)
}
