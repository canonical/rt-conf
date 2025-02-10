package ui

import (
	cmp "github.com/canonical/rt-conf/src/ui/components"
)

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
