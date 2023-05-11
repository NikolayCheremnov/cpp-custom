package ll1

import (
	"cpp-custom/internal/datatype"
	"cpp-custom/internal/il"
	"cpp-custom/internal/polis"
	"cpp-custom/internal/stree"
	"cpp-custom/logger"
	"errors"
	"strconv"
	"strings"
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
	START_EXPRESSION                 = "#se"
	END_EXPRESSION_WITH_DECLARATION  = "#eewd"
	// END_EXPRESSION                              = "#ee"
	START_ASSIGNMENT                            = "#sa"
	END_EXPRESSION_WITH_ASSIGNMENT              = "#eewa"
	END_EXPRESSIONT_WITH_ASSUGNMENT_FOR         = "#eewaf"
	END_EXPRESSION_WITH_CONDITION               = "#eewc"
	END_EXPRESSION_WITH_ASSIGNMENT_AT_FOR_INIT  = "#eewafi"
	END_EXPRESSION_WITH_DECLARATION_AT_FOR_INIT = "#eewdfi"
	DEFAULT_INITIALIZATION                      = "#di"
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
		return sfl(root, operations, ctx)
	case END_FOR_LOOP:
		return efl(root, operations, ctx)
	case FOR_LOOP_COUNTER_DECLARATION:
		return flcd(root, ctx)
	case START_PROC_CALL:
		return spc(root, operations, ctx)
	case PROC_VARIABLE_ARG_TRANSFER:
		return pvat(root, operations, ctx)
	case PROC_CONSTANT_ARG_TRANSFER:
		return pcat(operations, ctx)
	case END_PROC_CALL:
		return epc(operations)
	case START_ASSIGNMENT:
		return sa(ctx)
	case START_EXPRESSION:
		return se(ctx)
	case END_EXPRESSION_WITH_DECLARATION:
		return eewd(root, operations, ctx)
	case END_EXPRESSION_WITH_ASSIGNMENT:
		return eewa(root, operations, ctx)
	case END_EXPRESSIONT_WITH_ASSUGNMENT_FOR:
		return eewaf(root, operations, ctx)
	case END_EXPRESSION_WITH_CONDITION:
		return eewc(root, operations, ctx)
	case END_EXPRESSION_WITH_ASSIGNMENT_AT_FOR_INIT:
		return eewafi(root, operations, ctx)
	case END_EXPRESSION_WITH_DECLARATION_AT_FOR_INIT:
		return eewdfi(root, operations, ctx)
	case DEFAULT_INITIALIZATION:
		return di(root, operations, ctx)
	//case END_EXPRESSION:
	//	return ee(root, operations, ctx)
	default:
		return errors.New("Bad operational symbol received: " + operational)
	}
}

func di(root *stree.Root, operations *il.Intermediate, ctx *context) error {
	defaultValue := datatype.NewIntDataValue(0)
	operations.Operations = append(operations.Operations, il.Operation{
		IlInstruction: il.CAST,
		LeftOperand: &il.Operand{
			Type:         il.CONSTANT,
			Identity:     strconv.Itoa(int(defaultValue.DataAsInt)),
			OperandValue: &defaultValue,
		},
		RightOperand: &il.Operand{
			Type:         il.TYPE,
			Identity:     root.CurrentNode.Variable.Type.FullName,
			OperandValue: nil,
		},
		Result: nil,
	})
	operations.Operations = append(operations.Operations, il.Operation{
		IlInstruction: il.STOR,
		LeftOperand: &il.Operand{
			Type:         il.VARIABLE,
			Identity:     root.CurrentNode.IlIdentifier,
			OperandValue: nil,
		},
		RightOperand: &il.Operand{
			Type:         il.OPERATION,
			Identity:     strconv.Itoa(len(operations.Operations) - 1),
			OperandValue: nil,
		},
		Result: nil,
	})
	return nil
}

func sa(ctx *context) error {
	ctx.assigmentTarget = ctx.identityLexeme
	return nil
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
		NodeType:     stree.VARIABLE,
		Identifier:   ctx.identityLexeme,
		IlIdentifier: ctx.identityLexeme + ctx.counterAsString(),
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
	if err := operations.ExtractArgumentFromStack(argNode.IlIdentifier); err != nil {
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

func sfl(root *stree.Root, operations *il.Intermediate, ctx *context) error {
	root.CurrentNode.Left = &stree.Node{
		NodeType:   stree.FOR_OPERATOR,
		Identifier: "",
		Variable:   nil,
		Parent:     root.CurrentNode,
		Left:       nil,
		Right:      nil,
	}
	root.CurrentNode = root.CurrentNode.Left
	// add pass operation to start cycle
	operations.Operations = append(operations.Operations, il.Operation{
		IlInstruction: il.PASS,
		LeftOperand:   nil,
		RightOperand:  nil,
		Result:        nil,
	})
	// add deffered jump oeration
	ctx.deferredOperations = append(ctx.deferredOperations, il.Operation{
		IlInstruction: il.JMP,
		LeftOperand: &il.Operand{
			Type:         il.OPERATION,
			Identity:     strconv.Itoa(len(operations.Operations) - 1),
			OperandValue: nil,
		},
		RightOperand: nil,
		Result:       nil,
	})
	return nil
}

func efl(root *stree.Root, operations *il.Intermediate, ctx *context) error {
	var deferredCycleOperations []il.Operation
	i := len(ctx.deferredOperations) - 1
	for {
		deferredCycleOperations = append(deferredCycleOperations, ctx.deferredOperations[i])
		if ctx.deferredOperations[i].IlInstruction == il.JMP {
			break
		} else {
			i--
		}
	}
	ctx.deferredOperations = ctx.deferredOperations[:i] // remove from deferred operations in context
	moveOperationsToGlobalIlCode(operations, deferredCycleOperations[:len(deferredCycleOperations)-1])
	jmpOperation := deferredCycleOperations[len(deferredCycleOperations)-1]
	jmpTarget, _ := strconv.Atoi(jmpOperation.LeftOperand.Identity)
	operations.Operations = append(operations.Operations, jmpOperation) // last jump without offset
	// add pass at end of cycle
	operations.Operations = append(operations.Operations, il.Operation{
		IlInstruction: il.PASS,
		LeftOperand:   nil,
		RightOperand:  nil,
		Result:        nil,
	})
	// find nearest JMPF with empty left operand but to jmpTarget
	for i := len(operations.Operations) - 1; i > jmpTarget; i-- {
		if operations.Operations[i].IlInstruction == il.JMPF && operations.Operations[i].LeftOperand == nil {
			operations.Operations[i].LeftOperand = &il.Operand{
				Type:         il.OPERATION,
				Identity:     strconv.Itoa(len(operations.Operations) - 1),
				OperandValue: nil,
			}
		}
	}
	// change tree at the end
	return comeBackToSubRoot(root, stree.FOR_OPERATOR)
}

func flcd(root *stree.Root, ctx *context) error {
	root.CurrentNode.Right = &stree.Node{
		NodeType:     stree.VARIABLE,
		Identifier:   ctx.identityLexeme,
		IlIdentifier: ctx.identityLexeme + ctx.counterAsString(),
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
			LeftOperand: &il.Operand{
				Type:         il.OPERATION,
				Identity:     strconv.Itoa(len(operations.Operations) - 2), // - 1 because CALL move to end
				OperandValue: nil,
			},
			RightOperand: nil,
			Result:       nil,
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
		Type:         il.VARIABLE,
		Identity:     variableNode.IlIdentifier,
		OperandValue: nil,
	}
	//operations.Operations[lastOperationWithNotFilledLeftOperandIndex+1].LeftOperand = &il.Operand{
	//	Type:         il.VARIABLE,
	//	Identity:     variableNode.Identifier,
	//	OperandValue: nil,
	//}
	ctx.release()
	return nil
}

func pcat(operations *il.Intermediate, ctx *context) error {
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
	intValue, _ := strconv.Atoi(ctx.constantLexeme)
	value := datatype.NewIntDataValue(int64(intValue))
	operations.Operations[lastOperationWithNotFilledLeftOperandIndex].LeftOperand = &il.Operand{
		Type:         il.CONSTANT,
		Identity:     ctx.constantLexeme,
		OperandValue: &value,
	}
	//  intValue, _ := strconv.Atoi(ctx.constantLexeme)
	//  value := datatype.NewIntDataValue(int64(intValue))
	//operations.Operations[lastOperationWithNotFilledLeftOperandIndex+1].LeftOperand = &il.Operand{
	//	Type:         il.CONSTANT,
	//	Identity:     ctx.constantLexeme,
	//	OperandValue: &value,
	//}
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

func se(ctx *context) error {
	ctx.expressionTokens = []string{}
	ctx.expressionTokens = append(ctx.expressionTokens, ctx.lastWrittenLex)
	ctx.isExpressionParsing = true
	return nil
}

func eewd(root *stree.Root, operations *il.Intermediate, ctx *context) error {
	return processStorExpressionResult(root.CurrentNode, root, operations, ctx, false)
}

func eewa(root *stree.Root, operations *il.Intermediate, ctx *context) error {
	storTarget := root.CurrentNode.FindUpByIdentifier(ctx.assigmentTarget)
	if storTarget == nil || storTarget.NodeType != stree.VARIABLE || !storTarget.Variable.IsMutable {
		return errors.New("Not found mutable variable '" + ctx.assigmentTarget + "'")
	}
	return processStorExpressionResult(storTarget, root, operations, ctx, false)
}

func eewaf(root *stree.Root, operations *il.Intermediate, ctx *context) error {
	storTarget := root.CurrentNode.FindUpByIdentifier(ctx.assigmentTarget)
	if storTarget == nil || storTarget.NodeType != stree.VARIABLE || !storTarget.Variable.IsMutable {
		return errors.New("Not found mutable variable '" + ctx.assigmentTarget + "'")
	}
	return processStorExpressionResult(storTarget, root, operations, ctx, true)
}

func eewc(root *stree.Root, operations *il.Intermediate, ctx *context) error {
	ctx.isExpressionParsing = false
	logger.Log(OperationsL, "completed expression: "+strings.Join(ctx.expressionTokens, " "))
	var leftOperand *il.Operand
	var expressionOperations []il.Operation = nil
	if len(ctx.expressionTokens) == 1 {
		// then only one operand in expression
		var err error
		leftOperand, err = extractSingleExpressionOperandFromCtx(root, ctx)
		if err != nil {
			return err
		}
	} else {
		// else need parse and process expression
		var err error
		expressionOperations, err = expressionTokensToOperationsSlice(root, ctx)
		if err != nil {
			return err
		}
		// need to insert it into il operations
		moveOperationsToGlobalIlCode(operations, expressionOperations)
		leftOperand = &il.Operand{
			Type:         il.OPERATION,
			Identity:     strconv.Itoa(len(operations.Operations) - 1),
			OperandValue: nil,
		}
	}
	// add type converting to bool
	operations.Operations = append(operations.Operations, il.Operation{
		IlInstruction: il.CAST,
		LeftOperand:   leftOperand,
		RightOperand: &il.Operand{
			Type:         il.TYPE,
			Identity:     "bool",
			OperandValue: nil,
		},
		Result: nil,
	})
	operations.Operations = append(operations.Operations, il.Operation{
		IlInstruction: il.JMPF,
		LeftOperand:   nil,
		RightOperand: &il.Operand{
			Type:         il.OPERATION,
			Identity:     strconv.Itoa(len(operations.Operations) - 1),
			OperandValue: nil,
		},
		Result: nil,
	})
	return nil
}

func eewdfi(root *stree.Root, operations *il.Intermediate, ctx *context) error {
	if err := eewd(root, operations, ctx); err != nil {
		return err
	}
	return upPassInCycleInit(operations, ctx)
}

func eewafi(root *stree.Root, operations *il.Intermediate, ctx *context) error {
	if err := eewa(root, operations, ctx); err != nil {
		return err
	}
	return upPassInCycleInit(operations, ctx)
}

//func ee(root *stree.Root, operations *il.Intermediate, ctx *context) error {
//	// parse expression
//	ctx.isExpressionParsing = false
//	logger.Log(OperationsL, "completed expression: "+strings.Join(ctx.expressionTokens, " "))
//	tokensWithUnaryMinus := polis.AddUnaryMinus(ctx.expressionTokens)
//	logger.Log(OperationsL, "completed expression with unary minus: "+strings.Join(tokensWithUnaryMinus, " "))
//	polisTokens := polis.ConvertToRPN(tokensWithUnaryMinus)
//	logger.Log(OperationsL, "polis expression: "+strings.Join(polisTokens, " "))
//	// generate triads for expression
//	var operationsStack []il.Operation
//	var operationsOut []il.Operation
//	// pass by tokens in reversed order
//	for i := len(polisTokens) - 1; i >= 0; i-- {
//		token := polisTokens[i]
//		if polis.IsOperator(token) {
//			// add new Operation to stack top with empty operands
//			operationsStack = append(operationsStack, il.Operation{
//				IlInstruction: il.GetOperatorByLexeme(token),
//				LeftOperand:   nil,
//				RightOperand:  nil,
//				Result:        nil,
//			})
//		} else { // else variable or constant
//			var operand il.Operand
//			if val, err := strconv.Atoi(token); err == nil {
//				// then constant
//				value := datatype.NewIntDataValue(int64(val))
//				operand = il.Operand{
//					Type:         il.CONSTANT,
//					Identity:     token,
//					OperandValue: &value,
//				}
//			} else {
//				// then variable
//				// need check in tree
//				variableNode := root.CurrentNode.FindUpByIdentifier(token)
//				if variableNode == nil || variableNode.NodeType != stree.VARIABLE {
//					return errors.New("identifier '" + token + "' not found")
//				}
//				operand = il.Operand{
//					Type:         il.VARIABLE,
//					Identity:     token,
//					OperandValue: nil,
//				}
//			}
//			// write operand to up operation (start from right)
//			if operationsStack[len(operationsStack)-1].RightOperand == nil {
//				operationsStack[len(operationsStack)-1].RightOperand = &operand
//			} else if operationsStack[len(operationsStack)-1].LeftOperand == nil {
//				operationsStack[len(operationsStack)-1].LeftOperand = &operand
//			} else {
//				panic("no empty operands places in operation")
//			}
//		}
//		// condense operations
//		for len(operationsStack) > 0 &&
//			operationsStack[len(operationsStack)-1].RightOperand != nil &&
//			operationsStack[len(operationsStack)-1].LeftOperand != nil {
//			// move full filled operation to out
//			operationsOut = append(operationsOut, operationsStack[len(operationsStack)-1])
//			// remove it from stack
//			operationsStack = operationsStack[:len(operationsStack)-1]
//			// write full filled operation as operand for another operation
//			if len(operationsStack) > 0 {
//				operationAsOperand := il.Operand{
//					Type:         il.OPERATION,
//					Identity:     strconv.Itoa(len(operationsOut) - 1),
//					OperandValue: nil,
//				}
//				if operationsStack[len(operationsStack)-1].RightOperand == nil {
//					operationsStack[len(operationsStack)-1].RightOperand = &operationAsOperand
//				} else if operationsStack[len(operationsStack)-1].LeftOperand == nil {
//					operationsStack[len(operationsStack)-1].LeftOperand = &operationAsOperand
//				} else {
//					panic("no empty operands places for operation operand in operation")
//				}
//			}
//		}
//	}
//	// move operations to global il-code
//	newStartIndex := len(operations.Operations)
//	for _, movedOperation := range operationsOut {
//		if movedOperation.LeftOperand.Type == il.OPERATION {
//			oldIdentityValue, _ := strconv.Atoi(movedOperation.LeftOperand.Identity)
//			movedOperation.LeftOperand.Identity = strconv.Itoa(newStartIndex + oldIdentityValue)
//		}
//		if movedOperation.RightOperand.Type == il.OPERATION {
//			oldIdentityValue, _ := strconv.Atoi(movedOperation.RightOperand.Identity)
//			movedOperation.RightOperand.Identity = strconv.Itoa(newStartIndex + oldIdentityValue)
//		}
//		operations.Operations = append(operations.Operations, movedOperation)
//	}
//	//// try to view calculated operations
//	//out := strings.Builder{}
//	//for i, o := range operationsOut {
//	//	out.WriteString(strconv.Itoa(i) + ") " + o.OperationAsString() + "\n")
//	//}
//	//logger.Log(OperationsL, "operations:\n"+out.String())
//	return nil
//}

// common functions

// 1. variable declaration
func variableDeclaration(root *stree.Root, ctx *context, isFirst, isMutable bool) error {
	if root.CurrentNode.FindUpByIdentifierToNodeType(ctx.identityLexeme, stree.COMPOSITE_OPERATOR) != nil {
		return errors.New("redeclaration of '" + ctx.identityLexeme + "'")
	}
	variableNode := &stree.Node{
		NodeType:     stree.VARIABLE,
		Identifier:   ctx.identityLexeme,
		IlIdentifier: ctx.identityLexeme + ctx.counterAsString(),
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

// gen single expression operand
func extractSingleExpressionOperandFromCtx(root *stree.Root, ctx *context) (*il.Operand, error) {
	if len(ctx.expressionTokens) != 1 {
		return nil, errors.New("can`t get single expression operand from tokens, len != 1")
	}
	token := ctx.expressionTokens[0]
	tokenNode := root.CurrentNode.FindUpByIdentifier(token)
	if val, err := strconv.Atoi(token); err == nil {
		value := datatype.NewIntDataValue(int64(val))
		// then constant
		return &il.Operand{
			Type:         il.CONSTANT,
			Identity:     token,
			OperandValue: &value,
		}, nil
	} else {
		// then variable
		variableNode := root.CurrentNode.FindUpByIdentifier(token)
		if variableNode == nil || variableNode.NodeType != stree.VARIABLE {
			return nil, errors.New("identifier '" + token + "' not found")
		}
		return &il.Operand{
			Type:         il.VARIABLE,
			Identity:     tokenNode.IlIdentifier,
			OperandValue: nil,
		}, nil
	}
}

// parse expression tokens to operations
func expressionTokensToOperationsSlice(root *stree.Root, ctx *context) ([]il.Operation,
	error) {
	tokensWithUnaryMinus := polis.AddUnaryMinus(ctx.expressionTokens)
	logger.Log(OperationsL, "completed expression with unary minus: "+strings.Join(tokensWithUnaryMinus, " "))
	polisTokens := polis.ConvertToRPN(tokensWithUnaryMinus)
	logger.Log(OperationsL, "polis expression: "+strings.Join(polisTokens, " "))
	// generate triads for expression
	var operationsStack []il.Operation
	var operationsOut []il.Operation
	// pass by tokens in reversed order
	for i := len(polisTokens) - 1; i >= 0; i-- {
		token := polisTokens[i]
		if polis.IsOperator(token) {
			// add new Operation to stack top with empty operands
			operationsStack = append(operationsStack, il.Operation{
				IlInstruction: il.GetOperatorByLexeme(token),
				LeftOperand:   nil,
				RightOperand:  nil,
				Result:        nil,
			})
		} else { // else variable or constant
			var operand il.Operand
			if val, err := strconv.Atoi(token); err == nil {
				// then constant
				value := datatype.NewIntDataValue(int64(val))
				operand = il.Operand{
					Type:         il.CONSTANT,
					Identity:     token,
					OperandValue: &value,
				}
			} else {
				// then variable
				// need check in tree
				variableNode := root.CurrentNode.FindUpByIdentifier(token)
				if variableNode == nil || variableNode.NodeType != stree.VARIABLE {
					return operationsOut, errors.New("identifier '" + token + "' not found")
				}
				operand = il.Operand{
					Type:         il.VARIABLE,
					Identity:     variableNode.IlIdentifier,
					OperandValue: nil,
				}
			}
			// write operand to up operation (start from right)
			if operationsStack[len(operationsStack)-1].RightOperand == nil {
				operationsStack[len(operationsStack)-1].RightOperand = &operand
			} else if operationsStack[len(operationsStack)-1].LeftOperand == nil {
				operationsStack[len(operationsStack)-1].LeftOperand = &operand
			} else {
				panic("no empty operands places in operation")
			}
		}
		// condense operations
		for len(operationsStack) > 0 &&
			operationsStack[len(operationsStack)-1].RightOperand != nil &&
			operationsStack[len(operationsStack)-1].LeftOperand != nil {
			// move full filled operation to out
			operationsOut = append(operationsOut, operationsStack[len(operationsStack)-1])
			// remove it from stack
			operationsStack = operationsStack[:len(operationsStack)-1]
			// write full filled operation as operand for another operation
			if len(operationsStack) > 0 {
				operationAsOperand := il.Operand{
					Type:         il.OPERATION,
					Identity:     strconv.Itoa(len(operationsOut) - 1),
					OperandValue: nil,
				}
				if operationsStack[len(operationsStack)-1].RightOperand == nil {
					operationsStack[len(operationsStack)-1].RightOperand = &operationAsOperand
				} else if operationsStack[len(operationsStack)-1].LeftOperand == nil {
					operationsStack[len(operationsStack)-1].LeftOperand = &operationAsOperand
				} else {
					panic("no empty operands places for operation operand in operation")
				}
			}
		}
	}
	return operationsOut, nil
}

func moveOperationsToGlobalIlCode(operations *il.Intermediate, operationsToMove []il.Operation) {
	newStartIndex := len(operations.Operations)
	for _, movedOperation := range operationsToMove {
		if movedOperation.LeftOperand != nil && movedOperation.LeftOperand.Type == il.OPERATION {
			oldIdentityValue, _ := strconv.Atoi(movedOperation.LeftOperand.Identity)
			movedOperation.LeftOperand.Identity = strconv.Itoa(newStartIndex + oldIdentityValue)
		}
		if movedOperation.RightOperand != nil && movedOperation.RightOperand.Type == il.OPERATION {
			oldIdentityValue, _ := strconv.Atoi(movedOperation.RightOperand.Identity)
			movedOperation.RightOperand.Identity = strconv.Itoa(newStartIndex + oldIdentityValue)
		}
		operations.Operations = append(operations.Operations, movedOperation)
	}
}

func processStorExpressionResult(storTarget *stree.Node, root *stree.Root, operations *il.Intermediate, ctx *context,
	needMoveToDeferredOperations bool) error {
	// end parse expression
	ctx.isExpressionParsing = false
	logger.Log(OperationsL, "completed expression: "+strings.Join(ctx.expressionTokens, " "))
	var rightOperand *il.Operand
	var expressionOperations []il.Operation = nil
	if len(ctx.expressionTokens) == 1 {
		// then only one operand in expression
		var err error
		rightOperand, err = extractSingleExpressionOperandFromCtx(root, ctx)
		if err != nil {
			return err
		}
	} else {
		// else need parse and process expression
		var err error
		expressionOperations, err = expressionTokensToOperationsSlice(root, ctx)
		if err != nil {
			return err
		}
		// need to insert it into il operations
		var identity string
		if !needMoveToDeferredOperations {
			moveOperationsToGlobalIlCode(operations, expressionOperations)
			identity = strconv.Itoa(len(operations.Operations) - 1)
		} else {
			identity = strconv.Itoa(len(expressionOperations) - 1)
		}
		rightOperand = &il.Operand{
			Type:         il.OPERATION,
			Identity:     identity,
			OperandValue: nil,
		}
	}
	// then we know right operand of expression and can add operation for stor without mutable checking
	if storTarget.NodeType != stree.VARIABLE {
		return errors.New("no variable for initialization in current node")
	}
	// add type converting and stor
	var storRightOperandIdentity string
	if !needMoveToDeferredOperations {
		storRightOperandIdentity = strconv.Itoa(len(operations.Operations))
	} else {
		storRightOperandIdentity = strconv.Itoa(len(expressionOperations))
	}

	castOperation := il.Operation{
		IlInstruction: il.CAST,
		LeftOperand:   rightOperand, // this is funny moment because rightOperand for stor, but it is left for cast
		RightOperand: &il.Operand{
			Type:         il.TYPE,
			Identity:     storTarget.Variable.Type.FullName,
			OperandValue: nil,
		},
		Result: nil,
	}
	storOperation := il.Operation{
		IlInstruction: il.STOR,
		LeftOperand: &il.Operand{
			Type:         il.VARIABLE,
			Identity:     storTarget.IlIdentifier,
			OperandValue: nil,
		},
		RightOperand: &il.Operand{
			Type:         il.OPERATION,
			Identity:     storRightOperandIdentity,
			OperandValue: nil,
		},
		Result: nil,
	}
	if !needMoveToDeferredOperations {
		operations.Operations = append(operations.Operations, castOperation)
		operations.Operations = append(operations.Operations, storOperation)
	} else {
		ctx.deferredOperations = append(ctx.deferredOperations, storOperation)
		ctx.deferredOperations = append(ctx.deferredOperations, castOperation)
		// add expression operations to deffered operations stack in reversed order
		if expressionOperations != nil {
			for i := len(expressionOperations) - 1; i >= 0; i-- {
				ctx.deferredOperations = append(ctx.deferredOperations, expressionOperations[i])
			}
		}
	}

	return nil
}

func upPassInCycleInit(operations *il.Intermediate, ctx *context) error {
	// find nearest PASS-operator
	passIndex := -1
	for i := len(operations.Operations) - 1; i >= 0; i-- {
		if operations.Operations[i].IlInstruction == il.PASS {
			passIndex = i
		}
	}
	if passIndex == -1 {
		return errors.New("no suitable pass operator")
	}
	// move it to current end
	oldPassIndex := passIndex
	for i := oldPassIndex + 1; i < len(operations.Operations); i++ {
		if operations.Operations[i].LeftOperand != nil && operations.Operations[i].LeftOperand.Type == il.OPERATION {
			oldIndex, _ := strconv.Atoi(operations.Operations[i].LeftOperand.Identity)
			operations.Operations[i].LeftOperand.Identity = strconv.Itoa(oldIndex - 1)
		}
		if operations.Operations[i].RightOperand != nil && operations.Operations[i].RightOperand.Type == il.OPERATION {
			oldIndex, _ := strconv.Atoi(operations.Operations[i].RightOperand.Identity)
			operations.Operations[i].RightOperand.Identity = strconv.Itoa(oldIndex - 1)
		}
	}
	operations.MoveOperationToEnd(passIndex)
	newPassIndex := len(operations.Operations) - 1
	// update jump instruction
	for i := len(ctx.deferredOperations) - 1; i >= 0; i-- {
		if ctx.deferredOperations[i].LeftOperand.Identity == strconv.Itoa(oldPassIndex) {
			ctx.deferredOperations[i].LeftOperand.Identity = strconv.Itoa(newPassIndex)
			break
		}
	}
	return nil
}
