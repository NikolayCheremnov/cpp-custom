package semanthoid

import (
	"cpp-custom/logger"
	"errors"
)

// procedures tree node
type ProcNode struct {
	// procedure prototype data
	Identifier        string
	ParamsCount       int
	ParamsTypesLabels []int
	ParamsIdentifiers []string
	// procedure position in source file
	ProcTextPos int
	ProcLine    int
	ProcLinePos int
	// links
	Next *ProcNode
}

var ProcRoot *ProcNode
var CurrentProc *ProcNode

func AddProcedureDescription(identifier string, paramsCount int, paramsTypesLabels []int, paramsIdentifiers []string,
	procTextPos int, procLine int, procLinePos int) error {
	node := &ProcNode{
		Identifier:        identifier,
		ParamsCount:       paramsCount,
		ParamsTypesLabels: paramsTypesLabels,
		ParamsIdentifiers: paramsIdentifiers,
		ProcTextPos:       procTextPos,
		ProcLine:          procLine,
		ProcLinePos:       procLinePos,
	}
	if ProcRoot == nil {
		ProcRoot = node
		CurrentProc = node
	} else {
		redefinition := FindFromRoot(identifier) // check redefinitions
		if redefinition != nil {
			return errors.New("the '" + identifier + "' procedure has already been defined")
		}
		CurrentProc.Next = node
		CurrentProc = node
	}
	logger.Log("procedures_tree_l", "create procedure '"+identifier+"' description\n"+ProcListToString())
	return nil
}

// find methods

func FindFromRoot(identifier string) *ProcNode {
	if ProcRoot == nil {
		return nil
	}
	return ProcRoot.Find(identifier)
}

func (node *ProcNode) Find(identifier string) *ProcNode {
	if node.Identifier == identifier {
		return node
	}
	if node.Next == nil {
		return nil
	}
	return node.Next.Find(identifier)
}
