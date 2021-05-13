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
var ExecutableProcedure *ProcNode

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

// exec procedure

// load procedure parameters into tree, set right branch direction for first data in procedure
// return proc node, error
func LoadProcedure(identifier string, paramsValues []*DataTypeValue) (*ProcNode, error) {
	proc := FindFromRoot(identifier)
	if proc == nil {
		return nil, errors.New("procedure '" + identifier + "' has no definition")
	}
	BranchDirection = "right"
	for paramIndex, paramIdentifier := range proc.ParamsIdentifiers {
		LoadProcedureParameter(paramIdentifier, proc.ParamsTypesLabels[paramIndex], paramsValues[paramIndex])
		BranchDirection = "left"
	}
	logger.Log("tree_l", "'"+identifier+"' context is loaded into tree\n"+TreeToString())
	if proc.ParamsCount > 0 {
		logger.Log("memory_l", "memory allocation for '"+identifier+"' parameters")
	}
	return proc, nil
}

func LoadProcedureParameter(paramIdentifier string, paramType int, paramValue *DataTypeValue) {
	node := &Node{
		NodeTypeLabel: Variable,
		Identifier:    paramIdentifier,
		DataTypeLabel: paramType,
		DataValue:     paramValue,
	}
	if Root == nil {
		Root = node
		Current = node
	} else {
		if BranchDirection == "right" {
			Current.Right = node
		} else {
			Current.Left = node
		}
		node.Parent = Current
		Current = node
	}
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
