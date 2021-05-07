package semanthoid

import (
	"errors"
)

func InsertionPoint() string {
	if Current.NodeTypeLabel == CompositeOperator && Current.Right == nil {
		return "right"
	}
	return "left"
}

func CreateDescription(descriptionType int, identifier string, dataType int, value *DataTypeValue) (*Node, error) {
	if descriptionType != VariableNode && descriptionType != ConstantNode {
		return nil, errors.New("unknown description type")
	}
	if dataType != IntType && dataType != BoolType {
		return nil, errors.New("unknown data type")
	}
	if value == nil {
		value = GetDefaultDataValue()
	}
	return &Node{
		NodeTypeLabel: VariableNode,
		Identifier:    identifier,
		DataTypeLabel: dataType,
		DataValue:     value,
	}, nil
}
