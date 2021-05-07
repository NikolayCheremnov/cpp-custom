package semanthoid

import (
	"errors"
)

var ProceduresRoot *Node
var LastProcedure *Node

// returns procedure node and error
func CreateProcedureDescription(identifier string, paramsCount int, paramsTypes []int, paramsIdentifiers []string, procedureBody *Node) error {
	redefinition := FindFromNodeAmongLeft(ProceduresRoot, ProcedureDescription, identifier)
	if redefinition != nil {
		return errors.New("there is already a procedure '" + identifier + "'")
	}
	if identifier == "main" && paramsCount > 0 {
		return errors.New("main must not contain parameters")
	}
	// procedure description subtree initialization in insertion
	node := Node{
		NodeTypeLabel: ProcedureDescription,
		Identifier:    identifier,
		ParamsCount:   paramsCount,
		ParamsTypes:   paramsTypes,
	}
	// set links: between current and procedure and between procedure and composite operator
	if ProceduresRoot == nil { // if first node in tree
		ProceduresRoot = &node
		LastProcedure = &node
	} else if LastProcedure.Left != nil {
		return errors.New("fatal error in semantic tree: not empty left branch")
	} else { // procedures descriptions is only left children or root
		LastProcedure.Left = &node
		node.Parent = LastProcedure
	}
	if paramsCount > 0 { // if procedure contains parameters
		for i, param := range paramsIdentifiers {
			for j, otherParams := range paramsIdentifiers {
				if i != j && param == otherParams { // param redefinition
					return errors.New("there is already a formal parameter '" + param + "'")
				}
			}
		}
		node.Right = &Node{
			NodeTypeLabel: Variable,
			Identifier:    paramsIdentifiers[0],
			DataTypeLabel: paramsTypes[0],
			DataValue:     GetDefaultDataValue(),
			Parent:        &node,
		}
		subCurrent := node.Right
		for i := 1; i < paramsCount; i++ {
			subCurrent.Left = &Node{
				NodeTypeLabel: Variable,
				Identifier:    paramsIdentifiers[i],
				DataTypeLabel: paramsTypes[i],
				DataValue:     GetDefaultDataValue(),
				Parent:        subCurrent,
			}
			subCurrent = subCurrent.Left
		}
		subCurrent.Left = procedureBody
		procedureBody.Parent = subCurrent
	} else { // if immediately the body
		node.Right = procedureBody
		procedureBody.Parent = &node
	}
	return nil
}
