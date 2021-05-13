package semanthoid

import (
	"cpp-custom/logger"
	"errors"
)

// procedure context
type ProcContext struct {
	Proc        *ProcNode
	ContextRoot *Node
	ContextLeaf *Node
}

// procedures stack

type ProcStack []*ProcContext

func (s *ProcStack) PushBack(context *ProcContext) {
	*s = append(*s, context)
}

func (s *ProcStack) PopBack() *ProcContext {
	if s.IsEmpty() {
		return nil
	}
	context := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return context
}

func (s *ProcStack) IsEmpty() bool {
	if len(*s) == 0 {
		return true
	}
	return false
}

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
var Stack ProcStack
var ExecutableProc *ProcNode

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
		redefinition := FindFromHead(identifier) // check redefinitions
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

// load procedure parameters into tree, create fork for first data in procedure
// return proc node, error
func LoadProcedure(identifier string, paramsValues []int) (*ProcNode, error) {
	// find proc from proc list
	proc := FindFromHead(identifier)
	if proc == nil {
		return nil, errors.New("procedure '" + identifier + "' has no definition")
	}
	if (proc.ParamsCount != 0 && paramsValues == nil) || proc.ParamsCount != len(paramsValues) {
		return nil, errors.New("invalid number of parameters for '" + identifier + "'")
	}
	// create fork for procedure
	CreateFork()
	// entering parameters in the tree
	for paramIndex, paramIdentifier := range proc.ParamsIdentifiers {
		var paramValue *DataTypeValue
		if proc.ParamsTypesLabels[paramIndex] == IntType {
			paramValue = &DataTypeValue{paramsValues[paramIndex], IntToBool(paramsValues[paramIndex])}
		} else {
			paramValue = &DataTypeValue{IntToBool(paramsValues[paramIndex]), IntToBool(paramsValues[paramIndex])}
		}
		LoadProcedureParameter(paramIdentifier, proc.ParamsTypesLabels[paramIndex], paramValue)
	}
	logger.Log("tree_l", "'"+identifier+"' context is loaded into tree\n"+TreeToString())
	if proc.ParamsCount > 0 {
		logger.Log("memory_l", "memory allocation for '"+identifier+"' parameters")
	}
	ExecutableProc = proc
	return proc, nil
}

func LoadProcedureParameter(paramIdentifier string, paramType int, paramValue *DataTypeValue) {
	node := &Node{
		NodeTypeLabel: Variable,
		Identifier:    paramIdentifier,
		DataTypeLabel: paramType,
		DataValue:     paramValue,
	}
	Current.Left = node
	node.Parent = Current
	Current = node
}

// unloading procedure context from stack to tree
func LoadProcedureContext() error {
	context := Stack.PopBack()
	if context == nil {
		return errors.New("attempt to extract procedure context from an empty stack")
	}
	if Root == nil {
		Root = context.ContextRoot
		Current = context.ContextLeaf
	} else {
		Current.Right = context.ContextRoot
		Current.Right.Parent = Current
		Current = context.ContextLeaf
	}
	logger.Log("stack_l", "'"+context.Proc.Identifier+"' moved from stack\n"+Stack.ToString())
	return nil
}

// saves the context of the executable procedure in stack
func SaveProcedureContext() error {
	context, err := RemoveProcedureContextFromTree()
	if err != nil {
		return err
	}
	Stack.PushBack(context)
	logger.Log("stack_l", "'"+context.Proc.Identifier+"' moved to stack\n"+Stack.ToString())
	return nil
}

func RemoveProcedureContextFromTree() (*ProcContext, error) {
	procSubTree := findNearestForkFromRootAmongLeft()
	if procSubTree == nil {
		return nil, errors.New("no executable procedures")
	}
	// context saving
	context := &ProcContext{
		Proc:        ExecutableProc,
		ContextRoot: procSubTree,
		ContextLeaf: Current,
	}
	// remove context from tree
	if procSubTree.Parent == nil {
		Root = nil
		Current = nil
	} else {
		Current = procSubTree.Parent
		Current.Right = nil
	}
	return context, nil
}

// find methods

func FindFromHead(identifier string) *ProcNode {
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
