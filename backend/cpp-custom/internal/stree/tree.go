package stree

import "strings"

// Root of the tree - contains global tree root and procedures list
type Root struct {
	RootNode    *Node
	CurrentNode *Node
}

func NewRoot() *Root {
	root := new(Root)
	root.CurrentNode = &Node{NodeType: ROOT, Variable: nil, Parent: nil, Left: nil, Right: nil}
	root.RootNode = root.CurrentNode
	return root
}

func (root *Root) AsString() string {
	builder := new(strings.Builder)
	root.RootNode.SubtreeAsString(builder, 0)
	return builder.String()
}
