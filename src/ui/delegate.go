package ui

import (
	"github.com/canonical/rt-conf/src/ui/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

func newItemDelegateMainMenu(keys *selectKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.Styles.SelectedDesc = styles.SelectedDesc
	d.Styles.SelectedTitle = styles.SelectedTitle

	d.Styles.NormalDesc = styles.NormalDesc
	d.Styles.NormalTitle = styles.NormalTitle

	d.Styles.DimmedDesc = styles.DimmedDesc
	d.Styles.DimmedTitle = styles.DimmedTitle

	d.Styles.FilterMatch = styles.FilterMatch

	help := []key.Binding{keys.choose}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

type selectKeyMap struct {
	choose key.Binding
	remove key.Binding
}

// Additional short help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d selectKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.choose,
	}
}

// Additional full help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d selectKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.choose,
		},
	}
}

func newDelegateKeyMapMainMenu() *selectKeyMap {
	return &selectKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
	}
}
