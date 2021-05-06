package semanthoid

import "errors"

func InsertionPoint() string {
	if Current.NodeTypeLabel == CompositeOperator && Current.Right == nil {
		return "right"
	}
	return "left"
}

func CreateDescription(descriptionType int, identifier string, dataType int, value *DataTypeValue) error {
	if descriptionType != VariableNode && descriptionType != ConstantNode {
		return errors.New("unknown description type")
	}
	if dataType != IntType && dataType != BoolType {
		return errors.New("unknown data type")
	}
	if value == nil {
		value = &DataTypeValue{0, 0} // default values
	}
	node := Node{
		NodeTypeLabel: VariableNode,
		Identifier:    identifier,
		DataTypeLabel: dataType,
		DataValue:     value,
	}
	if Current == nil { // if first node in tree
		Current = &node
		Root = Current
	} else if InsertionPoint() == "left" {
		redefinition := Current.FindUpInCurrentRightSubTree(VariableNode, identifier)
		if redefinition != nil {
			return errors.New("there is already a variable named '" + identifier + "'")
		}
		redefinition = Current.FindUpInCurrentRightSubTree(ConstantNode, identifier)
		if redefinition != nil {
			return errors.New("there is already a constant named '" + identifier + "'")
		}
		Current.Left = &node
		node.Parent = Current
		Current = &node
	} else { // if first in right subtree
		Current.Right = &node
		node.Parent = Current
		Current = &node
	}
	return nil
}
