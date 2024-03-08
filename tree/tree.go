package tree

import (
	"fmt"
	"strings"
)

const (
	elemPrefix     = `├─ `
	lastElemPrefix = `└─ `
	pipe           = `│  `
	padding        = " "
	newLine        = "\n"
)

type TreeNode interface { //nolint:revive // keeping the names as it is an interface and Node is a struct
	AddChild(TreeNode)
	Children() []TreeNode
	AddChildren([]TreeNode)
	Format(string, string) string
	Value() any
	SetValue(any)
	Level() int
	SetLevel(int)
	StringPadding() int
	SetStringPadding(int)
}

type Node struct {
	level         int
	parent        TreeNode
	children      []TreeNode
	value         any
	stringPadding int
}

func NewNode(v any) TreeNode {
	return &Node{
		value: v,
	}
}

func (n *Node) IsRoot() bool {
	return n.parent == nil
}

func (n *Node) IsLeaf() bool {
	return n.children == nil || len(n.children) == 0
}

func (n *Node) AddChild(child TreeNode) {
	if child == nil {
		return
	}
	child.SetLevel(n.level + 1)
	child.SetStringPadding(n.stringPadding)
	if n.children == nil {
		n.children = []TreeNode{child}
		return
	}
	n.children = append(n.children, child)
}

func (n *Node) Children() []TreeNode {
	return n.children
}

func (n *Node) AddChildren(children []TreeNode) {
	if children == nil {
		return
	}
	for _, c := range children {
		n.AddChild(c)
	}
}

func (n *Node) Format(parentPrefix, prefix string) string {
	var sb strings.Builder
	nodeLine := fmt.Sprintf("%s%s%v", parentPrefix, prefix, n.Value())
	sb.WriteString(nodeLine)
	childrenLen := len(n.Children())
	for i, c := range n.Children() {
		prnPfx := fmt.Sprintf("%s%s", parentPrefix, strings.Repeat(padding, n.StringPadding()))
		if c.Level() > 1 && prefix != lastElemPrefix {
			prnPfx = fmt.Sprintf("%s%s", parentPrefix, pipe)
		}
		if c.Level() == 1 {
			prnPfx = fmt.Sprintf("%s%s", parentPrefix, padding)
		}
		pfx := elemPrefix
		if i == (childrenLen - 1) {
			pfx = lastElemPrefix
		}
		sb.WriteString(newLine)
		c.SetStringPadding(n.StringPadding())
		sb.WriteString(c.Format(prnPfx, pfx))
	}
	return sb.String()
}

func (n *Node) SetValue(v any) {
	n.value = v
}

func (n *Node) Value() any {
	return n.value
}

func (n *Node) Level() int {
	return n.level
}

func (n *Node) SetLevel(l int) {
	n.level = l
	for _, c := range n.children {
		c.SetLevel(l + 1)
	}
}

func (n *Node) SetStringPadding(p int) {
	n.stringPadding = p
	for _, c := range n.children {
		c.SetStringPadding(p)
	}
}

func (n *Node) StringPadding() int {
	return n.stringPadding
}
