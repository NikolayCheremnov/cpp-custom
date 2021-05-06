package semanthoid

import (
	"errors"
	"strconv"
	"strings"
)

var Root *Node = nil
var Current *Node = nil

var DEBUG_MODE bool = false

type DataTypeValue struct {
	DataAsInt  int
	DataAsBool int
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

// display methods

func (node *Node) ToString() string {
	res := "node type " + strconv.Itoa(node.NodeTypeLabel) + ": identifier " + node.Identifier
	switch node.NodeTypeLabel { // special view for each node type
	case ProcedureDescription:
		res += ": params count " + strconv.Itoa(node.ParamsCount)
		if node.ParamsCount > 0 {
			res += ": params types"
			for _, paramType := range node.ParamsTypes {
				res += " " + strconv.Itoa(paramType)
			}
		}
		break
	case VariableNode, ConstantNode:
		res += ": data type " + strconv.Itoa(node.DataTypeLabel)
		res += ": data value bool " + strconv.Itoa(node.DataValue.DataAsBool) + " int " + strconv.Itoa(node.DataValue.DataAsInt)
		break
	default:
		res += ": there are no additional information for node"
	}
	return res
}

func TreeToString() string {
	if Root != nil {
		return Root.TreeToString(1)
	}
	return "{ nil }"
}

// ROOT-RIGHT-LEFT scheme
func (node *Node) TreeToString(offset int) string {
	lowerOffsetStr := strings.Repeat("\t", offset-1)
	upperOffsetStr := strings.Repeat("\t", offset)
	nodeStr := upperOffsetStr + node.ToString() + "\n"
	var left string
	if node.Left != nil {
		left = node.Left.TreeToString(offset + 1)
	} else {
		left = upperOffsetStr + "{ nil }\n"
	}
	var right string
	if node.Right != nil {
		right = node.Right.TreeToString(offset + 1)
	} else {
		right = upperOffsetStr + "{ nil }\n"
	}
	return lowerOffsetStr + "{\n" + nodeStr + right + left + lowerOffsetStr + "}\n"
}
