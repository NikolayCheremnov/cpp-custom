package ll1

import (
	"cpp-custom/internal/datatype"
	"cpp-custom/internal/stree"
	"cpp-custom/logger"
	"errors"
)

// operational symbols
const (
	START_PROC_DECLARATION            = "#spd"
	PROC_IDENTIFIER                   = "#pi"
	PROC_ARG_DECLARATION              = "#pad"
	END_PROC_DECLARATION              = "#epd"
	START_COMPOSITE_OP                = "#sco"
	END_COMPOSITE_OP                  = "#eco"
	START_VARIABLES_DECLARATION       = "#svd"
	END_VATIABLES_DECLARATION         = "#evd"
	FIRST_VARIABLE_DECLARATION        = "#fvd"
	VARIABLE_DECLARATION              = "#vd"
	START_NAMED_CONSTANTS_DECLARATION = "#sncd"
	END_NAMED_CONSTANTS_DECLARATION   = "#encd"
	FIRST_NAMED_CONSTANT_DECLARATION  = "#fncd"
	NAMED_CONSTANT_DECLARATION        = "#ncd"
	START_FOR_LOOP                    = "#sfl"
	END_FOR_LOOP                      = "#efl"
	FOR_LOOP_COUNTER_DECLARATION      = "#flcd"
)

func runOperational(operational string, root *stree.Root, ctx *context) error {
	switch operational {
	case START_PROC_DECLARATION:
		return spd(ctx)
	case PROC_IDENTIFIER:
		return pi(root, ctx)
	case PROC_ARG_DECLARATION:
		return pad(root, ctx)
	case END_PROC_DECLARATION:
		return epd(root)
	case START_COMPOSITE_OP:
		return sco(root)
	case END_COMPOSITE_OP:
		return eco(root)
	case START_VARIABLES_DECLARATION:
		return svd(ctx)
	case END_VATIABLES_DECLARATION:
		return evd(ctx)
	case FIRST_VARIABLE_DECLARATION:
		return fvd(root, ctx)
	case VARIABLE_DECLARATION:
		return vd(root, ctx)
	case START_NAMED_CONSTANTS_DECLARATION:
		return sncd(ctx)
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
	default:
		return errors.New("Bad operational symbol received: " + operational)
	}
}

// operational handlers
func spd(ctx *context) error {
	if err := ctx.setMode(PROCEDURE_DECLARATION_MODE); err != nil {
		return err
	}
	return nil
}

func pi(root *stree.Root, ctx *context) error {
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
	ctx.release()
	return nil
}

func pad(root *stree.Root, ctx *context) error {
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
	ctx.release()
	return nil
}

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

func svd(ctx *context) error {
	if err := ctx.setMode(VARIABLE_DECLARATION_MODE); err != nil {
		return err
	}
	ctx.identityLexeme = ""
	return nil
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

func sncd(ctx *context) error {
	if err := ctx.setMode(CONSTANT_DECLARATION_MODE); err != nil {
		return err
	}
	return nil
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
