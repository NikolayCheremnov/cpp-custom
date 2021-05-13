package semanthoid

import (
	"cpp-custom/logger"
	"errors"
)

func checkDescriptionConditions(identifier string, descriptionType int) error {
	if descriptionType != Variable && descriptionType != Constant {
		return errors.New("bad description type label")
	}
	if FindDataUpFromCurrentInCurrentRightSubTree(identifier) != nil {
		return errors.New("'" + identifier + "' already declared in this block")
	}
	return nil
}

func createDescriptionNode(descriptionType int, identifier string, dataType int, value *DataTypeValue) (*Node, error) {
	err := checkDescriptionConditions(identifier, descriptionType)
	if err != nil {
		return nil, err
	}
	return &Node{
		NodeTypeLabel: descriptionType,
		Identifier:    identifier,
		DataTypeLabel: dataType,
		DataValue:     value}, nil
}

func CreateGlobalDescription(descriptionType int, identifier string, dataType int, value *DataTypeValue) error {
	node, err := createDescriptionNode(descriptionType, identifier, dataType, value)
	if err != nil {
		return err
	}
	if Root == nil { // if first in tree
		Root = node
		Current = Root
	} else {
		Current.Left = node
		node.Parent = Current
		Current = node
	}
	logger.Log("memory_l", "memory allocation for global description '"+identifier+"'")
	logger.Log("tree_l", "created global description '"+identifier+"'\n"+TreeToString())
	return nil
}

func CreateLocalDescription(descriptionType int, identifier string, dataType int, value *DataTypeValue) error {
	node, err := createDescriptionNode(descriptionType, identifier, dataType, value)
	if err != nil {
		return err
	}
	Current.Left = node
	node.Parent = Current
	Current = node

	logger.Log("memory_l", "memory allocation for local description '"+identifier+"'")
	logger.Log("tree_l", "created local description '"+identifier+"'\n"+TreeToString())
	return nil
}
