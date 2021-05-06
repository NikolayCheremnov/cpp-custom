package parsenator

import (
	"cpp-custom/internal/lexinator"
	"cpp-custom/internal/semanthoid"
	"errors"
	"fmt"
	"io"
)

// analyzer structure
type Analyzer struct {
	scanner    lexinator.Scanner
	isPrepared bool
	// syntax Tree
	programTree *semanthoid.Node
	// the output stream of error messages
	writer io.Writer
}

// preparing the analyzer
func Preparing(srsFileName string, scannerErrWriter io.Writer, analyzerErrWriter io.Writer) (Analyzer, error) {
	// syntax analyzer
	A := Analyzer{programTree: nil, writer: analyzerErrWriter}
	// lexical analyzer
	scanner, err := lexinator.ScannerInitializing(srsFileName, scannerErrWriter)
	if err != nil {
		return A, err
	}
	A.scanner = scanner
	// symantic tree
	if semanthoid.Root != nil {
		// TODO: add saving parse tree
		semanthoid.Root = nil
		semanthoid.Current = nil
	}
	// ready for parsing
	A.isPrepared = true
	return A, nil
}

// the results of the error message
func (A *Analyzer) printPanicError(msg string) {
	textPos, line, linePos := A.scanner.StorePosValues()
	_, err := fmt.Fprintf(A.writer, "error: %s position: %d line: %d line position: %d\n", msg, textPos, line, linePos)
	if err != nil {
		panic(err)
	}
	// maybe temporarily
	panic(errors.New("completed with an error. see the error logs")) // critical completion
}

func (A *Analyzer) printError(msg string) {
	textPos, line, linePos := A.scanner.StorePosValues()
	_, err := fmt.Fprintf(A.writer, "error: %s position: %d line: %d line position: %d\n", msg, textPos, line, linePos)
	if err != nil {
		panic(err)
	}
}

// handlers for the nonterminals

// <глобальные описания> -> e или ((<описание процедуры> | <описание> | ; |) + <глобальные описания>)
func (A *Analyzer) GlobalDescriptions() error {
	if !A.isPrepared {
		return errors.New("can't start the analysis: the analyzer is not prepared")
	}

	textPos, line, linePos := A.scanner.StorePosValues()
	lexType, lex := A.scanner.Scan()

	for lexType != lexinator.End {
		if lexType == lexinator.Void { // <описание процедуры>
			A.scanner.RestorePosValues(textPos, line, linePos)
			A.procedureDescription()
		} else if lexType == lexinator.Long ||
			lexType == lexinator.Short ||
			lexType == lexinator.Int ||
			lexType == lexinator.Bool ||
			lexType == lexinator.Const { // <описание>
			A.scanner.RestorePosValues(textPos, line, linePos)
			A.description()
		} else if lexType != lexinator.Semicolon { // then must be ';'
			A.printPanicError("invalid lexeme '" + lex + "', expected ';'")
		}
		textPos, line, linePos = A.scanner.StorePosValues()
		lexType, lex = A.scanner.Scan()
	}
	A.printError("there are no syntax level errors")
	return nil
}

// <описание параметров> ->
// return paramsCount and paramsTypes
func (A *Analyzer) parameterDescription() (int, []int) {
	var textPos, line, linePos int
	var lexType int
	var lex string
	isFirst, paramsCount := true, 0
	var paramsTypes []int
	for isFirst || lexType == lexinator.Comma {
		isFirst = false
		paramsTypes = append(paramsTypes, A._type())
		paramsCount++
		lexType, lex = A.scanner.Scan()
		if lexType != lexinator.Id {
			A.printPanicError("'" + lex + "' is not an identifier")
		}
		textPos, line, linePos = A.scanner.StorePosValues()
		lexType, lex = A.scanner.Scan()
	}
	A.scanner.RestorePosValues(textPos, line, linePos)
	return paramsCount, paramsTypes
}

// <параметры> -> идентификатор | константа U , <параметры> | e
// return params count
func (A *Analyzer) parameters() int {
	var paramsCount int = 0
	isFirst := true
	var textPos, line, linePos int
	var lexType int
	var lex string
	for isFirst || lexType == lexinator.Comma {
		isFirst = false
		lexType, lex = A.scanner.Scan()
		if lexType != lexinator.Id && lexType != lexinator.IntConst {
			A.printPanicError("'" + lex + "' is not an identifier or constant")
		}
		paramsCount++
		textPos, line, linePos = A.scanner.StorePosValues()
		lexType, lex = A.scanner.Scan()
	}
	A.scanner.RestorePosValues(textPos, line, linePos)
	return paramsCount
}

// <эл.выр.> -> (<выражение>) | идентификатор | константа
func (A *Analyzer) simplestExpr() {
	lexType, lex := A.scanner.Scan()
	if lexType == lexinator.OpeningBracket { // ( <выражение> )
		A.expression()
		lexType, lex = A.scanner.Scan()
		if lexType != lexinator.ClosingBracket {
			A.printPanicError("invalid lexeme '" + lex + "', expected ')'")
		}
	} else if lexType != lexinator.Id && lexType != lexinator.IntConst {
		A.printPanicError("'" + lex + "' not allowed in the expression")
	}
}

// <множитель> -> <эл.выр.> U e | * U / U % <эл.выр.>
func (A *Analyzer) multiplier() {
	A.simplestExpr() // <эл.выр.>
	textPos, line, linePos := A.scanner.StorePosValues()
	lexType, _ := A.scanner.Scan()
	for lexType == lexinator.Mul || lexType == lexinator.Div || lexType == lexinator.Mod {
		A.simplestExpr()
		textPos, line, linePos = A.scanner.StorePosValues()
		lexType, _ = A.scanner.Scan()
	}
	A.scanner.RestorePosValues(textPos, line, linePos)
}

// <процедура> -> идентификатор ( ) | идентификатор ( <параметры> )
func (A *Analyzer) procedure() {
	lexType, lex := A.scanner.Scan()
	if lexType != lexinator.Id {
		A.printPanicError("'" + lex + "' is not an identifier")
	}
	procedureDescription := semanthoid.FindFromNodeAmongLeft(semanthoid.Root, semanthoid.ProcedureDescription, lex)
	paramsCount := 0
	if procedureDescription == nil {
		A.printPanicError("undefined procedure '" + lex + "'")
	}
	lexType, lex = A.scanner.Scan()
	if lexType != lexinator.OpeningBracket {
		A.printPanicError("invalid lexeme '" + lex + "', expected '('")
	}
	textPos, line, linePos := A.scanner.StorePosValues()
	lexType, lex = A.scanner.Scan()
	if lexType != lexinator.ClosingBracket {
		A.scanner.RestorePosValues(textPos, line, linePos)
		paramsCount = A.parameters()
		lexType, lex = A.scanner.Scan()
		if lexType != lexinator.ClosingBracket {
			A.printPanicError("invalid lexeme '" + lex + "', expected ')'")
		}
	}
	if paramsCount != procedureDescription.ParamsCount {
		A.printPanicError("invalid arguments number in '" + procedureDescription.Identifier + "' calling")
	}
}

// <слагаемое> -> <множитель> U +- | e
func (A *Analyzer) summand() {
	A.multiplier() // <множитель>
	textPos, line, linePos := A.scanner.StorePosValues()
	lexType, _ := A.scanner.Scan()
	for lexType == lexinator.Plus || lexType == lexinator.Minus {
		A.multiplier()
		textPos, line, linePos = A.scanner.StorePosValues()
		lexType, _ = A.scanner.Scan()
	}
	A.scanner.RestorePosValues(textPos, line, linePos)
}

// <выражение> -> + | - | e U <слагаемое> + +- == <= >= < > <слагаемое> | e
func (A *Analyzer) expression() *semanthoid.DataTypeValue {
	textPos, line, linePos := A.scanner.StorePosValues()
	lexType, _ := A.scanner.Scan()
	if lexType != lexinator.Plus && lexType != lexinator.Minus {
		A.scanner.RestorePosValues(textPos, line, linePos)
	}
	A.summand() // <слагаемое>
	textPos, line, linePos = A.scanner.StorePosValues()
	lexType, _ = A.scanner.Scan()
	for lexType == lexinator.Plus || lexType == lexinator.Minus ||
		lexType == lexinator.Equ || lexType == lexinator.LessEqu || lexType == lexinator.MoreEqu ||
		lexType == lexinator.Less || lexType == lexinator.More {
		A.multiplier()
		textPos, line, linePos = A.scanner.StorePosValues()
		lexType, _ = A.scanner.Scan()
	}
	A.scanner.RestorePosValues(textPos, line, linePos)
	return &semanthoid.DataTypeValue{0, 0} // TODO: add expression value calculation
}

// <оператор for>
func (A *Analyzer) forOperator() {
	textPos, line, linePos := A.scanner.StorePosValues()
	lexType, lex := A.scanner.Scan()
	if lexType != lexinator.Semicolon { // if not empty
		if lexType != lexinator.Id { // if type
			A.scanner.RestorePosValues(textPos, line, linePos)
			A._type()
			lexType, lex = A.scanner.Scan()
			if lexType != lexinator.Id {
				A.printPanicError("'" + lex + "' is not an identifier")
			}
		}
		lexType, lex = A.scanner.Scan()
		if lexType != lexinator.Assignment {
			A.printPanicError("invalid lexeme '" + lex + "', expected '='")
		}
		A.expression()
	}
}

// <присваивание> -> идентификатор = <выражение>
// returns identifier and expression value
func (A *Analyzer) assigment() (identifier string, value *semanthoid.DataTypeValue) {
	lexType, lex := A.scanner.Scan()
	if lexType != lexinator.Id {
		A.printPanicError("'" + lex + "' is not an identifier")
	}
	identifier = lex
	lexType, lex = A.scanner.Scan()
	if lexType != lexinator.Assignment {
		A.printPanicError("invalid lexeme '" + lex + "', expected '='")
	}
	value = A.expression()
	return identifier, value
}

// <переменная> -> идентификатор U e | = <выражение>
func (A *Analyzer) variable() (string, *semanthoid.DataTypeValue) {
	lexType, lex := A.scanner.Scan()
	if lexType != lexinator.Id {
		A.printPanicError("'" + lex + "' is not an identifier")
	}
	identifier := lex // variable identifier
	var value *semanthoid.DataTypeValue = nil
	textPos, line, linePos := A.scanner.StorePosValues()
	lexType, lex = A.scanner.Scan()
	if lexType == lexinator.Assignment {
		A.expression() // TODO: add variable initializing with expression
	} else {
		A.scanner.RestorePosValues(textPos, line, linePos)
	}
	return identifier, value
}

// <for> -> for ( <оператор for> ; U <выражение> | e U ; U <присваивание> | e U ) <оператор>
func (A *Analyzer) _for() {
	lexType, lex := A.scanner.Scan()
	if lexType != lexinator.For {
		A.printPanicError("invalid lexeme '" + lex + "', expected 'for'")
	}
	lexType, lex = A.scanner.Scan()
	if lexType != lexinator.OpeningBracket {
		A.printPanicError("invalid lexeme '" + lex + "', expected '('")
	}
	A.forOperator()
	lexType, lex = A.scanner.Scan()
	if lexType != lexinator.Semicolon {
		A.printPanicError("invalid lexeme '" + lex + "', expected ';'")
	}
	textPos, line, linePos := A.scanner.StorePosValues()
	lexType, lex = A.scanner.Scan()
	if lexType != lexinator.Semicolon {
		A.scanner.RestorePosValues(textPos, line, linePos)
		A.expression()
		lexType, lex = A.scanner.Scan()
		if lexType != lexinator.Semicolon {
			A.printPanicError("invalid lexeme '" + lex + "', expected ';'")
		}
	}
	textPos, line, linePos = A.scanner.StorePosValues()
	lexType, lex = A.scanner.Scan()
	if lexType != lexinator.ClosingBracket {
		A.scanner.RestorePosValues(textPos, line, linePos)
		A.assigment()
		lexType, lex = A.scanner.Scan()
		if lexType != lexinator.ClosingBracket {
			A.printPanicError("invalid lexeme '" + lex + "', expected ';'")
		}
	}
	A.operator()
}

// <константы>
func (A *Analyzer) constants() {
	lexType, lex := A.scanner.Scan()
	if lexType != lexinator.Const {
		A.printPanicError("invalid lexeme '" + lex + "', expected 'const'")
	}
	constsType := A._type()
	var textPos, line, linePos int
	isFirst := true
	for isFirst || lexType == lexinator.Comma {
		isFirst = false
		identifier, value := A.assigment()
		err := semanthoid.CreateDescription(semanthoid.ConstantNode, identifier, constsType, value)
		if err != nil {
			A.printPanicError(err.Error())
		}
		textPos, line, linePos = A.scanner.StorePosValues()
		lexType, lex = A.scanner.Scan()
	}
	A.scanner.RestorePosValues(textPos, line, linePos)
}

// <переменные> -> long int | short int | int | bool U e | <присваивание> U e | , <присваивание>
func (A *Analyzer) variables() {
	varsType := A._type()
	var lexType, textPos, line, linePos int
	isFirst := true
	for isFirst || lexType == lexinator.Comma {
		isFirst = false
		identifier, value := A.variable()
		err := semanthoid.CreateDescription(semanthoid.VariableNode, identifier, varsType, value)
		if err != nil {
			A.printPanicError(err.Error())
		}
		textPos, line, linePos = A.scanner.StorePosValues()
		lexType, _ = A.scanner.Scan()
	}
	A.scanner.RestorePosValues(textPos, line, linePos)
}

// <оператор> -> <составной оператор> | <for> | <процедура> ; | <присваивание>; | ;
func (A *Analyzer) operator() {
	textPos, line, linePos := A.scanner.StorePosValues()
	lexType, lex := A.scanner.Scan()
	if lexType == lexinator.OpeningBrace { // составной оператор
		A.scanner.RestorePosValues(textPos, line, linePos)
		A.compositeOperator()
	} else if lexType == lexinator.For { // for
		A.scanner.RestorePosValues(textPos, line, linePos)
		A._for()
	} else if lexType == lexinator.Id { // процедура или присваивание
		lexType, lex = A.scanner.Scan()
		if lexType == lexinator.OpeningBracket { // процедура
			A.scanner.RestorePosValues(textPos, line, linePos)
			A.procedure()
		} else if lexType == lexinator.Assignment { // присваивание
			A.scanner.RestorePosValues(textPos, line, linePos)
			A.assigment()
		} else {
			A.printPanicError("'" + lex + "' is not an procedure or assigment")
		}
		lexType, lex = A.scanner.Scan()
		if lexType != lexinator.Semicolon {
			A.printPanicError("invalid lexeme '" + lex + "', expected ';'")
		}
	} else if lexType != lexinator.Semicolon {
		A.printPanicError("invalid lexeme '" + lex + "', expected ';'")
	}
}

// <описание процедуры> -> void идентификатор ( U <описание параметров> | e U ) <составной оператор>
func (A *Analyzer) procedureDescription() {
	// procedure node data
	var procedureIdentifier string
	var paramsCount int
	var paramsTypes []int
	// parsing
	lexType, lex := A.scanner.Scan()
	if lexType != lexinator.Void { // void
		A.printPanicError("invalid lexeme '" + lex + "', expected 'void'")
	}
	lexType, lex = A.scanner.Scan()
	if lexType != lexinator.Id { // identifier
		A.printPanicError("'" + lex + "' is not an identifier")
	}
	procedureIdentifier = lex
	lexType, lex = A.scanner.Scan()
	if lexType != lexinator.OpeningBracket { // (
		A.printPanicError("invalid lexeme '" + lex + "', expected '('")
	}
	textPos, line, linePos := A.scanner.StorePosValues()
	lexType, lex = A.scanner.Scan()
	if lexType == lexinator.Long || lexType == lexinator.Int ||
		lexType == lexinator.Short || lexType == lexinator.Bool { // <описание параметров>
		A.scanner.RestorePosValues(textPos, line, linePos)
		paramsCount, paramsTypes = A.parameterDescription()
		lexType, lex = A.scanner.Scan()
		if lexType != lexinator.ClosingBracket {
			A.printPanicError("invalid lexeme '" + lex + "', expected ')'")
		}
	} else if lexType != lexinator.ClosingBracket { // )
		A.printPanicError("invalid lexeme '" + lex + "', expected ')'")
	}
	procedureNode, err := semanthoid.CreateProcedureDescription(procedureIdentifier, paramsCount, paramsTypes)
	if err != nil {
		A.printPanicError(err.Error())
	}
	A.compositeOperator()
	semanthoid.Current = procedureNode // set procedure node as current
}

// <составной оператор> -> { <операторы и описания> }
func (A *Analyzer) compositeOperator() {
	compositeOperatorNode := semanthoid.CreateCompositeOperator() // store operator node
	lexType, lex := A.scanner.Scan()
	if lexType != lexinator.OpeningBrace {
		A.printPanicError("invalid lexeme '" + lex + "', expected '{'")
	}
	A.operatorsAndDescriptions()
	lexType, lex = A.scanner.Scan()
	if lexType != lexinator.ClosingBrace {
		A.printPanicError("invalid lexeme '" + lex + "', expected '}'")
	}
	semanthoid.Current = compositeOperatorNode // restore operator node
}

// <операторы и описания> -> e | <операторы> U e | <операторы и описания>  | <описания> + e | <операторы и описания>
func (A *Analyzer) operatorsAndDescriptions() {
	for {
		textPos, line, linePos := A.scanner.StorePosValues()
		lexType, _ := A.scanner.Scan()
		A.scanner.RestorePosValues(textPos, line, linePos)
		if lexType == lexinator.OpeningBrace || lexType == lexinator.For ||
			lexType == lexinator.Id || lexType == lexinator.Semicolon { // if operator
			A.operator()
		} else if lexType == lexinator.Long || lexType == lexinator.Short ||
			lexType == lexinator.Int || lexType == lexinator.Bool || lexType == lexinator.Const { // if description
			A.description()
		} else { // e
			break
		}
	}
}

// <описание> -> <переменные>; | <константы>;
func (A *Analyzer) description() {
	textPos, line, linePos := A.scanner.StorePosValues()
	lexType, lex := A.scanner.Scan()
	if lexType == lexinator.Long ||
		lexType == lexinator.Short ||
		lexType == lexinator.Int ||
		lexType == lexinator.Bool { // <переменные>
		A.scanner.RestorePosValues(textPos, line, linePos)
		A.variables()
	} else if lexType == lexinator.Const { // <константы>
		A.scanner.RestorePosValues(textPos, line, linePos)
		A.constants()
	} else {
		A.printPanicError("'" + lex + "'" + "does not name a type")
	}
	lexType, lex = A.scanner.Scan()
	if lexType != lexinator.Semicolon {
		A.printPanicError("invalid lexeme '" + lex + "', expected ';'")
	}
}

// <тип> -> long int | short int | int | bool
// return type label
func (A *Analyzer) _type() int {
	lexType, lex := A.scanner.Scan()
	if lexType != lexinator.Long && lexType != lexinator.Short &&
		lexType != lexinator.Int && lexType != lexinator.Bool {
		A.printPanicError("'" + lex + "'" + "does not name a type")
	} else if lexType == lexinator.Long || lexType == lexinator.Short {
		lexType, lex = A.scanner.Scan()
		if lexType != lexinator.Int {
			A.printPanicError("invalid lexeme '" + lex + "', expected 'int'")
		}
	}
	if lexType == lexinator.Bool {
		return semanthoid.IntType
	} else {
		return semanthoid.BoolType
	}
}
