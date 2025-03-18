package ui

import (
	"fmt"
	"log"
	"strings"

	"github.com/canonical/rt-conf/src/data"
	"github.com/canonical/rt-conf/src/kcmd"
	cmp "github.com/canonical/rt-conf/src/ui/components"
	"github.com/canonical/rt-conf/src/ui/config"
	"github.com/canonical/rt-conf/src/ui/styles"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

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
	iConf      data.InternalConfig
	// keys     *listKeyMap
}

type kcmdKeyMap struct {
	Up         key.Binding
	Down       key.Binding
	Left       key.Binding
	Right      key.Binding
	Select     key.Binding
	Help       key.Binding
	Quit       key.Binding
	CursorMode key.Binding
	Back       key.Binding
	Home       key.Binding
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k kcmdKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right, k.Select}, // first column
		{k.Help, k.Quit, k.Back, k.Home, k.Quit},  // second column
	}
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k kcmdKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Back, k.Quit}
}

func newkcmdMenuListKeyMap() *kcmdKeyMap {
	return &kcmdKeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "move left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "move left"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter", "space"),
			key.WithHelp("enter", "select menu"),
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

func newKcmdMenuModel(c *data.InternalConfig) KcmdlineMenuModel {
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

func (m KcmdlineMenuModel) Init() tea.Cmd { return textinput.Blink }

const (
	isolcpusIndex = iota
	nohzIndex
	nohzFullIndex
	kthreadsCPUsIndex
	irqaffinityIndex
	applyButtonIndex
	backButtonIndex
)

var placeholders_text = []string{
	config.CpuListPlaceholder,
	"Enter on or off",
	config.CpuListPlaceholder,
	config.CpuListPlaceholder,
	config.CpuListPlaceholder,
}

// TODO: fix the problem with the j,k keys being logged
func (m *KcmdlineMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	log.Println("---(kcmdlineMenuUpdate - start")
	cmds := make([]tea.Cmd, len(m.Inputs))

	totalIrqItems := len(m.Inputs) + 2
	index := cmp.NewNavigation(&m.FocusIndex, &totalIrqItems)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {

		case key.Matches(msg, m.keys.Select),
			key.Matches(msg, m.keys.Up),
			key.Matches(msg, m.keys.Down),
			key.Matches(msg, m.keys.Left),
			key.Matches(msg, m.keys.Right):

			// Handle navigation between the buttons
			if m.FocusIndex == applyButtonIndex &&
				key.Matches(msg, m.keys.Left) {
				index.Prev()
			}
			if m.FocusIndex == applyButtonIndex &&
				key.Matches(msg, m.keys.Right) {
				index.Next()
			}
			if m.FocusIndex == backButtonIndex &&
				key.Matches(msg, m.keys.Right) {
				index.Next()
			}
			if m.FocusIndex == backButtonIndex &&
				key.Matches(msg, m.keys.Left) {
				index.Prev()
			}

			// log.Println("focusIndex on Update: ", m.kcmdFocusIndex)
			// Validate the inputs

			valid := m.AreValidInputs()
			// log.Println("isValid: ", valid)

			// Handle [ Back ] button
			if m.FocusIndex == backButtonIndex &&
				key.Matches(msg, m.keys.Select) {
				// log.Println("pressed [ BACK ]: Back to main menu")
				m.Nav.PrevMenu()
			}

			// Did the user press enter while the apply button was focused?
			// TODO: improve mapping of len(m.inputs) to the apply button
			if key.Matches(msg, m.keys.Select) &&
				m.FocusIndex == len(m.Inputs) && valid {

				// log.Println("Apply changes")

				valid := m.AreValidInputs()

				if !valid {
					break
				}
				var empty int
				for i := range m.Inputs {
					v := m.Inputs[i].Value()
					if v == "" {
						empty++
					}
				}
				if empty == len(m.Inputs) {
					m.errorMsg = "\n\nAll fields are empty\n\n\n"
					break
				}

				m.iConf.Data.KernelCmdline.IsolCPUs = m.Inputs[isolcpusIndex].Value()

				m.iConf.Data.KernelCmdline.Nohz = m.Inputs[nohzIndex].Value()

				m.iConf.Data.KernelCmdline.NohzFull = m.Inputs[nohzFullIndex].Value()

				m.iConf.Data.KernelCmdline.KthreadCPUs = m.Inputs[kthreadsCPUsIndex].Value()

				m.iConf.Data.KernelCmdline.IRQaffinity = m.Inputs[irqaffinityIndex].Value()

				msgs, err := kcmd.ProcessKcmdArgs(&m.iConf)
				if err != nil {
					m.errorMsg = "Failed to process kernel cmdline args: " +
						err.Error()
					break
				}

				m.concl.logMsg = msgs
				m.concl.renderLog = true
				m.Nav.SetNewMenu(config.KCMD_CONCLUSION_VIEW_ID)

				// TODO: this needs to return a tea.Cmd (or maybe not)
				// TODO: Apply the changes call the kcmdline funcs
			}

			// Cycle indexes
			if key.Matches(msg, m.keys.Up) {
				index.Prev()
			}

			if key.Matches(msg, m.keys.Down) ||
				key.Matches(msg, m.keys.Select) {
				index.Next()
			}

			cmds := make([]tea.Cmd, len(m.Inputs))
			for i := 0; i <= len(m.Inputs)-1; i++ {
				if i == m.FocusIndex {
					// Set focused state
					cmds[i] = m.Inputs[i].Focus()
					m.Inputs[i].PromptStyle = styles.FocusedStyle
					m.Inputs[i].TextStyle = styles.FocusedStyle
					m.Inputs[i].Placeholder = placeholders_text[i]
					continue
				}
				// Remove focused state
				m.Inputs[i].Blur()
				m.Inputs[i].PromptStyle = styles.NoStyle
				m.Inputs[i].TextStyle = styles.NoStyle
				m.Inputs[i].Placeholder = ""
			}
		}
	}
	for i := range m.Inputs {
		m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
	}

	return m, tea.Batch(cmds...)
}

// TODO: fix the padding of the kcmline view
// * NOTE: in comparison with the main menu, the title is shifted to left
// * Not only the tittle but the hole view is shifted to the left

func (m KcmdlineMenuModel) View() string {
	var s string // the view

	title := styles.InnerMenuStyle("Configuring Kernel Cmdline Parameters")

	// The inputs
	var b strings.Builder
	for i := range m.Inputs {
		b.WriteString(m.Inputs[i].View())
		if i < len(m.Inputs)-1 {
			b.WriteRune('\n')
		}
	}

	apply_button := cmp.NewButton("Apply")
	back_button := cmp.NewButton("Back")
	apply_button.SetBlurred()
	back_button.SetBlurred()

	// TODO: add space between the [ Apply ] and [ Back ] buttons
	if m.FocusIndex == len(m.Inputs) {
		apply_button.SetFocused()
		back_button.SetBlurred()

	} else if m.FocusIndex == len(m.Inputs)+1 {
		apply_button.SetBlurred()
		back_button.SetFocused()
	}

	// [ Back ] [ Apply ] buttons
	fmt.Fprintf(&b, "\n\n%s\n\n",
		styles.JoinHorizontal(
			back_button.Render(),
			apply_button.Render(),
		))

	body := b.String()
	// TODO: Adding padding to the bottom and top of [body] and remove new lines

	//TODO: this needs to be dropped
	helpView := m.help.View(m.keys)

	// TODO: fix this mess
	height := (m.Height -
		strings.Count(title, "\n") -
		strings.Count(helpView, "\n") -
		strings.Count(m.errorMsg, "\n") -
		strings.Count(body, "\n") - 4) / 2 // TODO: fix those magic numbers

	// NOTE: *- 4 * because:
	// "\n\n" (jumps 2 lines) after the title
	// Before the line with the [ Back ] [ Apply ] buttons there are 2 lines

	// NOTE: * / 2 * (divide by two) because:
	// we want to add padding between to the top
	// and bottom of the help view

	if height < 0 {
		height = 1
	}

	s +=
		title +
			"\n\n" +
			body +
			strings.Repeat("\n", height) +
			"\n" +
			styles.ErrorMessageStyle(m.errorMsg) +
			strings.Repeat("\n", height) +
			helpView
	return styles.AppStyle.Render(s)
}
