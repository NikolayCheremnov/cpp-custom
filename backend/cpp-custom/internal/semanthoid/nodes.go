package semanthoid

var Root *Node = nil
var Current *Node = nil

var DEBUG_MODE = false

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

// find variable
func FindVariableUpFromCurrent(identifier string) *Node {
	if Current == nil {
		return nil
	}
	return Current.FindDataUp(identifier)
}

func (node *Node) FindVariableUp(identifier string) *Node {
	if node.NodeTypeLabel == Variable && node.Identifier == identifier {
		return node
	}
	if node.Parent == nil {
		return nil
	}
	return node.Parent.FindVariableUp(identifier)
}

// find data
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
