package stree

import (
	"cpp-custom/internal/datatype"
	"strings"
)

// NODE TYPES which starts right subtrees
const (
	ROOT               = "ROOT"
	VARIABLE           = "VARIABLE"
	PROCEDURE          = "PROCEDURE"
	COMPOSITE_OPERATOR = "COMPOSITE_OPERATOR"
	FOR_OPERATOR       = "FOR_OPERATOR"
)

type Node struct {
	NodeType     string
	Identifier   string
	IlIdentifier string
	Variable     *datatype.Variable
	// nodes pointers
	Parent *Node
	Right  *Node
	Left   *Node
}

func (node *Node) FindUpByNodeType(nodeType string) *Node {
	if node.NodeType == nodeType {
		return node
	}
	if node.Parent == nil {
		return nil
	}
	return node.Parent.FindUpByNodeType(nodeType)
}

func (node *Node) FindUpByIdentifierToNodeType(identifier string, stopNodeType string) *Node {
	if node.Identifier == identifier {
		return node
	}
	if node.Parent == nil || node.NodeType == stopNodeType {
		return nil
	}
	return node.Parent.FindUpByIdentifier(identifier)
}

func (node *Node) FindUpByIdentifier(identifier string) *Node {
	if node.Identifier == identifier {
		return node
	}
	if node.Parent == nil {
		return nil
	}
	return node.Parent.FindUpByIdentifier(identifier)
}

func (node *Node) NodeAsString() string {
	ilIdentifier := ""
	if node.IlIdentifier != "" {
		ilIdentifier = "(" + node.IlIdentifier + ")"
	}
	if node.Variable != nil {
		return node.NodeType + " : " + node.Identifier + ilIdentifier + " : " + node.Variable.VariableToString()
	}
	return node.NodeType + " : " + node.Identifier + ilIdentifier
}

func (node *Node) SubtreeAsString(builder *strings.Builder, offset int) {
	printOffset := func(o int) {
		for i := 0; i < o; i++ {
			builder.WriteString("\t")
		}
	}
	printOffset(offset)
	builder.WriteString("(" + node.NodeAsString() + "\n")
	if node.Left != nil {
		node.Left.SubtreeAsString(builder, offset+1)
	} else {
		printOffset(offset + 1)
		builder.WriteString("(nil)\n")
	}
	if node.Right != nil {
		node.Right.SubtreeAsString(builder, offset+1)
	} else {
		printOffset(offset + 1)
		builder.WriteString("(nil)\n")
	}
	printOffset(offset)
	builder.WriteString(")\n")
}
