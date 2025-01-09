package ui

import "errors"

type menuOpt int

const (
	mainMenu menuOpt = iota
	kcmdlineMenu
	irqAffinityMenu
	powerManagementMenu
)

// TODO: start to use these functions to improve Menu navigation
type MenuNode struct {
	Menu     menuOpt
	Children []*MenuNode
}

// MenuNavigator handles navigation within the menu tree
type MenuNavigator struct {
	CurrentNode *MenuNode   // Current position in the menu
	ParentStack []*MenuNode // Stack to keep track of parent nodes
}

// NewMenuNavigator initializes the navigator at the root of the menu tree
func NewMenuNavigator(root *MenuNode) *MenuNavigator {
	return &MenuNavigator{
		CurrentNode: root,
		ParentStack: []*MenuNode{},
	}
}

// AddChild adds a child node (submenu) to the current menu node
func (node *MenuNode) AddChild(child *MenuNode) {
	node.Children = append(node.Children, child)
}

// LinkedList represents the linked list
type LinkedList struct {
	Head *MenuNode
}

// Next navigates into a submenu by index
func (navigator *MenuNavigator) Next(index int) error {
	if index < 0 || index >= len(navigator.CurrentNode.Children) {
		return errors.New("invalid index: no such child menu")
	}

	// Push current node to the stack and move to the child
	navigator.ParentStack = append(navigator.ParentStack, navigator.CurrentNode)
	navigator.CurrentNode = navigator.CurrentNode.Children[index]
	return nil
}

// Previous navigates back to the parent menu
func (navigator *MenuNavigator) Previous() error {
	if len(navigator.ParentStack) == 0 {
		return errors.New("already at the top level; cannot go back")
	}

	// Pop the parent node from the stack
	navigator.CurrentNode = navigator.ParentStack[len(navigator.ParentStack)-1]
	navigator.ParentStack = navigator.ParentStack[:len(navigator.ParentStack)-1]
	return nil
}

// ******************** MENU NAVIGATION ********************

func (m *Model) NextIndex() {
	m.focusIndex++
	if m.focusIndex > len(m.inputs) {
		m.focusIndex = 0
	}
}

func (m *Model) PrevIndex() {
	m.focusIndex--
	if m.focusIndex == -1 {
		m.focusIndex = len(m.inputs)
	}
}
