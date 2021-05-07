package semanthoid

import (
	"errors"
)

var Root *Node = nil
var Current *Node = nil

var DEBUG_MODE bool = false

type DataTypeValue struct {
	DataAsInt  int
	DataAsBool int
}

func GetDefaultDataValue() *DataTypeValue {
	return &DataTypeValue{0, 0}
}

type Node struct {
	NodeTypeLabel int
	Identifier    string

	// for procedure
	ParamsCount int
	ParamsTypes []int

	// for variables
	DataTypeLabel int
	DataValue     *DataTypeValue

	// nodes pointers
	Right  *Node
	Left   *Node
	Parent *Node
}

func CreateTree(node *Node) error {
	if Current != nil {
		Current = node
		return nil
	}
	return errors.New("empty tree root")
}

// find procedures and methods
func FindDownLeft(node *Node, nodeType int, identifier string) *Node {
	if node == nil || (node.NodeTypeLabel == nodeType && node.Identifier == identifier) {
		return node
	}
	return FindDownLeft(node.Left, nodeType, identifier)
}

func FindFromNodeAmongLeft(node *Node, nodeType int, identifier string) *Node {
	if node == nil || (node.NodeTypeLabel == nodeType && node.Identifier == identifier) {
		return node
	}
	return FindFromNodeAmongLeft(node.Left, nodeType, identifier)
}

func (node *Node) FindUpInCurrentRightSubTree(nodeType int, identifier string) *Node {
	if node.Parent == nil || node.Parent.Right == node { // if node is root or root of right subtree
		if node.NodeTypeLabel == nodeType && node.Identifier == identifier {
			return node
		} else {
			return nil
		}
	} else if node.NodeTypeLabel == nodeType && node.Identifier == identifier {
		return node
	} else {
		return node.Parent.FindUpInCurrentRightSubTree(nodeType, identifier)
	}
}
