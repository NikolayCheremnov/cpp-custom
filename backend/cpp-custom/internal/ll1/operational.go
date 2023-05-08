package ll1

import (
	"cpp-custom/internal/datatype"
	"cpp-custom/internal/il"
	"cpp-custom/internal/stree"
	"cpp-custom/logger"
	"errors"
	"strconv"
)

// operational symbols
const (
	PROC_IDENTIFIER                  = "#pi"
	PROC_ARG_DECLARATION             = "#pad"
	END_PROC_DECLARATION             = "#epd"
	START_COMPOSITE_OP               = "#sco"
	END_COMPOSITE_OP                 = "#eco"
	END_VATIABLES_DECLARATION        = "#evd"
	FIRST_VARIABLE_DECLARATION       = "#fvd"
	VARIABLE_DECLARATION             = "#vd"
	END_NAMED_CONSTANTS_DECLARATION  = "#encd"
	FIRST_NAMED_CONSTANT_DECLARATION = "#fncd"
	NAMED_CONSTANT_DECLARATION       = "#ncd"
	START_FOR_LOOP                   = "#sfl"
	END_FOR_LOOP                     = "#efl"
	FOR_LOOP_COUNTER_DECLARATION     = "#flcd"
	START_PROC_CALL                  = "#spc"
	END_PROC_CALL                    = "#epc"
	PROC_VARIABLE_ARG_TRANSFER       = "#pvat"
	PROC_CONSTANT_ARG_TRANSFER       = "#pcat"
)

func runOperational(operational string, root *stree.Root, operations *il.Intermediate, ctx *context) error {
	switch operational {
	case PROC_IDENTIFIER:
		return pi(root, operations, ctx)
	case PROC_ARG_DECLARATION:
		return pad(root, operations, ctx)
	case END_PROC_DECLARATION:
		return epd(root)
	case START_COMPOSITE_OP:
		return sco(root)
	case END_COMPOSITE_OP:
		return eco(root)
	case END_VATIABLES_DECLARATION:
		return evd(ctx)
	case FIRST_VARIABLE_DECLARATION:
		return fvd(root, ctx)
	case VARIABLE_DECLARATION:
		return vd(root, ctx)
	case END_NAMED_CONSTANTS_DECLARATION:
		return encd(ctx)
	case FIRST_NAMED_CONSTANT_DECLARATION:
		return fncd(root, ctx)
	case NAMED_CONSTANT_DECLARATION:
		return ncd(root, ctx)
	case START_FOR_LOOP:
		return sfl(root)
	case END_FOR_LOOP:
		return efl(root)
	case FOR_LOOP_COUNTER_DECLARATION:
		return flcd(root, ctx)
	case START_PROC_CALL:
		return spc(root, operations, ctx)
	case PROC_VARIABLE_ARG_TRANSFER:
		return pvat(root, operations, ctx)
	case PROC_CONSTANT_ARG_TRANSFER:
		return pcat(root, operations, ctx)
	case END_PROC_CALL:
		return epc(operations)
	default:
		return errors.New("Bad operational symbol received: " + operational)
	}
}

func pi(root *stree.Root, operations *il.Intermediate, ctx *context) error {
	// check redeclaration
	if root.CurrentNode.FindUpByIdentifier(ctx.identityLexeme) != nil {
		return errors.New("redeclaration of '" + ctx.identityLexeme + "'")
	}
	// then add procedure node to tree
	root.CurrentNode.Left = &stree.Node{
		NodeType:   stree.PROCEDURE,
		Identifier: ctx.identityLexeme,
		Variable:   nil,
		Parent:     root.CurrentNode,
		Left:       nil,
		Right:      nil,
	}
	root.CurrentNode = root.CurrentNode.Left
	logger.Log(TreeL, "declare proc '"+ctx.identityLexeme+"'")
	// add operation declare proc
	if err := operations.DeclareProcedure(ctx.identityLexeme); err != nil {
		return err
	}
	logger.Log(OperationsL, "declare proc '"+ctx.identityLexeme+"'")
	ctx.release()
	return nil
}

func pad(root *stree.Root, operations *il.Intermediate, ctx *context) error {
	if root.CurrentNode.FindUpByIdentifierToNodeType(ctx.identityLexeme, stree.PROCEDURE) != nil {
		return errors.New("redeclaration of '" + ctx.identityLexeme + "'")
	}
	argNode := &stree.Node{
		NodeType:   stree.VARIABLE,
		Identifier: ctx.identityLexeme,
		Variable: &datatype.Variable{
			IsMutable: true,
			Type:      datatype.Type{FullName: ctx.getFullType()},
		},
		Parent: root.CurrentNode,
		Left:   nil,
		Right:  nil,
	}
	if root.CurrentNode.NodeType == stree.PROCEDURE {
		root.CurrentNode.Right = argNode
	} else {
		root.CurrentNode.Left = argNode
	}
	root.CurrentNode = argNode
	//
	if err := operations.ExtractArgumentFromStack(ctx.identityLexeme); err != nil {
		return err
	}
	//
	ctx.release()
	return nil
}

// operational symbols handlers

func epd(root *stree.Root) error {
	return comeBackToSubRoot(root, stree.PROCEDURE)
}

func sco(root *stree.Root) error {
	compositeOpNode := &stree.Node{
		NodeType:   stree.COMPOSITE_OPERATOR,
		Identifier: "",
		Variable:   nil,
		Parent:     root.CurrentNode,
		Left:       nil,
		Right:      nil,
	}
	if root.CurrentNode.NodeType == stree.PROCEDURE || root.CurrentNode.NodeType == stree.FOR_OPERATOR {
		// in procedure to a right subtree
		root.CurrentNode.Right = compositeOpNode
		root.CurrentNode = root.CurrentNode.Right
	} else {
		// args already in right subtree
		root.CurrentNode.Left = compositeOpNode
		root.CurrentNode = root.CurrentNode.Left
	}
	return nil
}

func eco(root *stree.Root) error {
	return comeBackToSubRoot(root, stree.COMPOSITE_OPERATOR)
}

func evd(ctx *context) error {
	ctx.release()
	return nil
}

func fvd(root *stree.Root, ctx *context) error {
	return variableDeclaration(root, ctx, true, true)
}

func vd(root *stree.Root, ctx *context) error {
	return variableDeclaration(root, ctx, false, true)
}

func encd(ctx *context) error {
	ctx.release()
	return nil
}

func fncd(root *stree.Root, ctx *context) error {
	return variableDeclaration(root, ctx, true, false)
}

func ncd(root *stree.Root, ctx *context) error {
	return variableDeclaration(root, ctx, false, false)
}

func sfl(root *stree.Root) error {
	root.CurrentNode.Left = &stree.Node{
		NodeType:   stree.FOR_OPERATOR,
		Identifier: "",
		Variable:   nil,
		Parent:     root.CurrentNode,
		Left:       nil,
		Right:      nil,
	}
	root.CurrentNode = root.CurrentNode.Left
	return nil
}

func efl(root *stree.Root) error {
	return comeBackToSubRoot(root, stree.FOR_OPERATOR)
}

func flcd(root *stree.Root, ctx *context) error {
	root.CurrentNode.Right = &stree.Node{
		NodeType:   stree.VARIABLE,
		Identifier: ctx.identityLexeme,
		Variable: &datatype.Variable{
			IsMutable: true,
			Type:      datatype.Type{FullName: ctx.getFullType()},
		},
		Parent: root.CurrentNode,
		Left:   nil,
		Right:  nil,
	}
	root.CurrentNode = root.CurrentNode.Right
	ctx.identityLexeme = ""
	ctx.typeSubLexeme = ""
	ctx.typeLexeme = ""
	return nil
}

func spc(root *stree.Root, operations *il.Intermediate, ctx *context) error {
	// 1. check that procedure existing
	procNode := root.CurrentNode.FindUpByIdentifier(ctx.identityLexeme)
	if procNode == nil {
		return errors.New("procedure '" + ctx.identityLexeme + "' has no declaration")
	}
	// 2. add proc calling instruction
	if err := operations.CallProcedure(procNode.Identifier); err != nil {
		return err
	}
	// 3. go for all args
	for argNode := procNode.Right; argNode.NodeType != stree.COMPOSITE_OPERATOR; argNode = argNode.Left {
		// add argument type cast operation (with empty left operand)
		operations.Operations = append(operations.Operations, il.Operation{
			IlInstruction: il.CAST,
			LeftOperand:   nil, // from is nil
			RightOperand: &il.Operand{
				Type:         il.TYPE,
				Identity:     argNode.Variable.Type.FullName,
				OperandValue: nil,
			},
			Result: nil,
		})
		// and push arg to stack (with empty left operand)
		operations.Operations = append(operations.Operations, il.Operation{
			IlInstruction: il.PUSH,
			LeftOperand:   nil,
			RightOperand:  nil,
			Result:        nil,
		})
	}
	ctx.release()
	return nil
}

func pvat(root *stree.Root, operations *il.Intermediate, ctx *context) error {
	variableNode := root.CurrentNode.FindUpByIdentifier(ctx.identityLexeme)
	if variableNode == nil || variableNode.NodeType != stree.VARIABLE {
		return errors.New("identifier '" + ctx.identityLexeme + "' not found")
	}
	// find operations for fill left operands
	lastProcCallOperationIndex := operations.FindLastProcCallOperationIndex()
	if lastProcCallOperationIndex == -1 {
		return errors.New("invalid operations structure")
	}
	// find last not filled operations
	lastOperationWithNotFilledLeftOperandIndex := -1
	for i := lastProcCallOperationIndex + 1; i < len(operations.Operations); i++ {
		if operations.Operations[i].LeftOperand == nil {
			lastOperationWithNotFilledLeftOperandIndex = i
			break
		}
	}
	if lastOperationWithNotFilledLeftOperandIndex == -1 {
		return errors.New("too many arguments for procedure '" +
			operations.Operations[lastProcCallOperationIndex].LeftOperand.Identity + "'")
	}
	// initialize left operations
	operations.Operations[lastOperationWithNotFilledLeftOperandIndex].LeftOperand = &il.Operand{
		Type:         il.TYPE,
		Identity:     variableNode.Variable.Type.FullName,
		OperandValue: nil,
	}
	operations.Operations[lastOperationWithNotFilledLeftOperandIndex+1].LeftOperand = &il.Operand{
		Type:         il.VARIABLE,
		Identity:     variableNode.Identifier,
		OperandValue: nil,
	}
	ctx.release()
	return nil
}

func pcat(root *stree.Root, operations *il.Intermediate, ctx *context) error {
	// find operations for fill left operands
	lastProcCallOperationIndex := operations.FindLastProcCallOperationIndex()
	if lastProcCallOperationIndex == -1 {
		return errors.New("invalid operations structure")
	}
	// find last not filled operations
	lastOperationWithNotFilledLeftOperandIndex := -1
	for i := lastProcCallOperationIndex + 1; i < len(operations.Operations); i++ {
		if operations.Operations[i].LeftOperand == nil {
			lastOperationWithNotFilledLeftOperandIndex = i
			break
		}
	}
	if lastOperationWithNotFilledLeftOperandIndex == -1 {
		return errors.New("too many arguments for procedure '" +
			operations.Operations[lastProcCallOperationIndex].LeftOperand.Identity + "'")
	}
	// initialize left operations
	operations.Operations[lastOperationWithNotFilledLeftOperandIndex].LeftOperand = &il.Operand{
		Type:         il.TYPE,
		Identity:     "int", // only int constants
		OperandValue: nil,
	}
	intValue, _ := strconv.Atoi(ctx.constantLexeme)
	value := datatype.NewIntDataValue(int64(intValue))
	operations.Operations[lastOperationWithNotFilledLeftOperandIndex+1].LeftOperand = &il.Operand{
		Type:         il.CONSTANT,
		Identity:     ctx.constantLexeme,
		OperandValue: &value,
	}
	ctx.release()
	return nil
}

func epc(operations *il.Intermediate) error {
	// find operations for fill left operands
	lastProcCallOperationIndex := operations.FindLastProcCallOperationIndex()
	if lastProcCallOperationIndex == -1 {
		return errors.New("invalid operations structure")
	}
	// find last not filled operations
	lastOperationWithNotFilledLeftOperandIndex := -1
	for i := lastProcCallOperationIndex + 1; i < len(operations.Operations); i++ {
		if operations.Operations[i].LeftOperand == nil {
			lastOperationWithNotFilledLeftOperandIndex = i
			break
		}
	}
	if lastOperationWithNotFilledLeftOperandIndex != -1 {
		return errors.New("too few arguments for procedure '" +
			operations.Operations[lastProcCallOperationIndex].LeftOperand.Identity + "'")
	}
	// move lastProcCallOperation to end of operations
	operations.MoveOperationToEnd(lastProcCallOperationIndex)
	return nil
}

// common functions

// 1. variable declaration
func variableDeclaration(root *stree.Root, ctx *context, isFirst, isMutable bool) error {
	if root.CurrentNode.FindUpByIdentifierToNodeType(ctx.identityLexeme, stree.COMPOSITE_OPERATOR) != nil {
		return errors.New("redeclaration of '" + ctx.identityLexeme + "'")
	}
	variableNode := &stree.Node{
		NodeType:   stree.VARIABLE,
		Identifier: ctx.identityLexeme,
		Variable: &datatype.Variable{
			IsMutable: isMutable,
			Type:      datatype.Type{FullName: ctx.getFullType()},
		},
		Parent: root.CurrentNode,
		Left:   nil,
		Right:  nil,
	}
	if isFirst && root.CurrentNode.NodeType == stree.COMPOSITE_OPERATOR {
		root.CurrentNode.Right = variableNode
		root.CurrentNode = root.CurrentNode.Right
	} else {
		root.CurrentNode.Left = variableNode
		root.CurrentNode = root.CurrentNode.Left
	}
	ctx.identityLexeme = ""
	return nil
}

// 2. come back to subroot
func comeBackToSubRoot(root *stree.Root, nodeType string) error {
	node := root.CurrentNode.FindUpByNodeType(nodeType)
	if node == nil {
		return errors.New("can`t find " + nodeType + " node")
	}
	root.CurrentNode = node
	return nil
}
