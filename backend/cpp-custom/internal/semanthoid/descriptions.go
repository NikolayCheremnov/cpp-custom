package semanthoid

import (
	"cpp-custom/logger"
	"errors"
)

func InsertionPoint() string {
	if Current.NodeTypeLabel == CompositeOperator && Current.Right == nil {
		return "right"
	}
	return "left"
}

func CreateGlobalDescription(descriptionType int, identifier string, dataType int, value *DataTypeValue) error {
	if descriptionType != Variable && descriptionType != Constant {
		return errors.New("bad description type label")
	}
	node := &Node{
		NodeTypeLabel: descriptionType,
		Identifier:    identifier,
		DataTypeLabel: dataType,
		DataValue:     value}
	if Root == nil { // if first in tree
		Root = node
		Current = Root
	} else {
		Current.Left = node
		node.Parent = Current
		Current = node
	}
	logger.Log("memory_l", "memory allocation for "+identifier+"\n"+TreeToString())
	return nil
}

func CreateDescription(descriptionType int, identifier string, dataType int, value *DataTypeValue) (*Node, error) {
	if descriptionType != Variable && descriptionType != Constant {
		return nil, errors.New("unknown description type")
	}
	if dataType != IntType && dataType != BoolType {
		return nil, errors.New("unknown data type")
	}
	if value == nil {
		value = GetDefaultDataValue()
	}
	return &Node{
		NodeTypeLabel: Variable,
		Identifier:    identifier,
		DataTypeLabel: dataType,
		DataValue:     value,
	}, nil
}
