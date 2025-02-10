package ui

import (
	"github.com/canonical/rt-conf/src/data"
	cmp "github.com/canonical/rt-conf/src/ui/components"
	"github.com/charmbracelet/bubbles/list"
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

type KcmdlineConclusion struct {
	Nav       *cmp.MenuNav // Menu Navigation instance
	keys      *kcmdKeyMap
	Width     int
	Height    int
	logMsg    []string
	renderLog bool
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
