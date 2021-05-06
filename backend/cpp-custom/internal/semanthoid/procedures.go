package semanthoid

import "errors"

// returns procedure node and error
func CreateProcedureDescription(identifier string, paramsCount int, paramsTypes []int) (*Node, error) {
	node := Node{
		NodeTypeLabel: ProcedureDescription,
		Identifier:    identifier,
		ParamsCount:   paramsCount,
		ParamsTypes:   paramsTypes,
	}
	// set links: between current and procedure and between procedure and composite operator
	if Current == nil { // if first node in tree
		Current = &node
		Root = &node
	} else if Current.Left != nil {
		return nil, errors.New("fatal error in semantic tree: not empty left branch")
	} else { // procedures descriptions is only left children or root
		redefinition := Current.FindUpInCurrentRightSubTree(ProcedureDescription, identifier)
		if redefinition != nil {
			return nil, errors.New("there is already a procedure '" + identifier + "'")
		}
		Current.Left = &node
	}
	node.Parent = Current
	node.Right = &Node{NodeTypeLabel: CompositeOperator, Identifier: CompositeOperatorIdentifier}
	node.Right.Parent = &node
	Current = node.Right
	return &node, nil
}
