package semanthoid

var Root *Node = nil
var Current *Node = nil
var BranchDirection = "left"

var DEBUG_MODE bool = false

// data types

type DataTypeValue struct {
	DataAsInt  int
	DataAsBool int
}

func GetDefaultDataValue() *DataTypeValue {
	return &DataTypeValue{0, 0}
}

func IntToBool(dataAsInt int) int {
	if dataAsInt != 0 {
		return 1
	}
	return 0
}

func GoBoolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

// nodes tree

type Node struct {
	NodeTypeLabel int
	Identifier    string

	// for variables and constants
	DataTypeLabel int
	DataValue     *DataTypeValue

	// nodes pointers
	Right  *Node
	Left   *Node
	Parent *Node
}

// find procedures and methods

// up from current in nearest right subtree
func FindDataUpFromCurrentInCurrentRightSubTree(identifier string) *Node {
	if Root == nil {
		return nil
	}
	return Current.FindDataUpInCurrentRightSubTree(identifier)
}

func (node *Node) FindDataUpInCurrentRightSubTree(identifier string) *Node {
	if (node.NodeTypeLabel == Variable || node.NodeTypeLabel == Constant) && node.Identifier == identifier {
		return node
	}
	if node.Parent == nil || node.Parent.Right == node {
		return nil
	}
	return node.Parent.FindDataUpInCurrentRightSubTree(identifier)
}

// up from current
func FindDataUpFromCurrent(identifier string) *Node {
	if Current == nil {
		return nil
	}
	return Current.FindDataUp(identifier)
}

func (node *Node) FindDataUp(identifier string) *Node {
	if (node.NodeTypeLabel == Variable || node.NodeTypeLabel == Constant) && node.Identifier == identifier {
		return node
	}
	if node.Parent == nil {
		return nil
	}
	return node.Parent.FindDataUp(identifier)
}

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
