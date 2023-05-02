package ll1

type context struct {
	identityLexeme string
	typeLexeme     string
	typeSubLexeme  string
}

func (ctx *context) saveIdentity(identity string) {
	ctx.identityLexeme = identity
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
}
