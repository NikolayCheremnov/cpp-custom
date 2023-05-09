package ll1

import "cpp-custom/internal/il"

type context struct {
	identityLexeme string
	constantLexeme string
	typeLexeme     string
	typeSubLexeme  string
	operatorLexeme string
	// expression data
	isExpressionParsing bool
	expressionTokens    []string
	// last added lex
	lastWrittenLex string
	//
	assigmentTarget string
	//
	deferredOperations []il.Operation
}

func (ctx *context) saveConstant(constantLex string) {
	ctx.constantLexeme = constantLex
	ctx.lastWrittenLex = constantLex
}

func (ctx *context) saveIdentity(identity string) {
	ctx.identityLexeme = identity
	ctx.lastWrittenLex = identity
}

func (ctx *context) saveOperator(operatorLex string) {
	ctx.operatorLexeme = operatorLex
	ctx.lastWrittenLex = operatorLex
}

func (ctx *context) saveType(typeLex string) {
	if (ctx.typeLexeme == "short" || ctx.typeLexeme == "long") && typeLex == "int" {
		ctx.typeSubLexeme = ctx.typeLexeme
	} else {
		ctx.typeSubLexeme = ""
	}
	ctx.typeLexeme = typeLex
}

func (ctx *context) getFullType() string {
	if ctx.typeSubLexeme == "" {
		return ctx.typeLexeme
	}
	return ctx.typeSubLexeme + " " + ctx.typeLexeme
}

func (ctx *context) release() {
	ctx.identityLexeme = ""
	ctx.typeLexeme = ""
	ctx.typeSubLexeme = ""
	ctx.constantLexeme = ""
}
