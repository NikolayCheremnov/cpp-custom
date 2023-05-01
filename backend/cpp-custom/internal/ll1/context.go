package ll1

import "errors"

// context mode constants
const (
	FREE = iota
	PROCEDURE_DECLARATION_MODE
	VARIABLE_DECLARATION_MODE
	CONSTANT_DECLARATION_MODE
)

type context struct {
	contextMode    int
	identityLexeme string
	typeLexeme     string
	typeSubLexeme  string
}

func (ctx *context) setMode(mode int) error {
	if ctx.contextMode != FREE {
		return errors.New("try set context to not free context")
	}
	ctx.contextMode = mode
	return nil
}

func (ctx *context) saveIdentity(identity string) error {
	if ctx.contextMode != FREE && ctx.identityLexeme != "" {
		return errors.New("try write second identity to not free context")
	}
	ctx.identityLexeme = identity
	return nil
}

func (ctx *context) saveType(typeLex string) error {
	if (ctx.typeLexeme == "short" || ctx.typeLexeme == "long") && typeLex == "int" {
		ctx.typeSubLexeme = ctx.typeLexeme
		ctx.typeLexeme = typeLex
		return nil
	} else if ctx.contextMode != FREE && ctx.typeLexeme != "" {
		return errors.New("try write second type to not free context")
	}
	ctx.typeLexeme = typeLex
	return nil
}

func (ctx *context) getFullType() string {
	if ctx.typeSubLexeme == "" {
		return ctx.typeLexeme
	}
	return ctx.typeSubLexeme + " " + ctx.typeLexeme
}

func (ctx *context) release() {
	ctx.contextMode = FREE
	ctx.identityLexeme = ""
	ctx.typeLexeme = ""
	ctx.typeSubLexeme = ""
}
