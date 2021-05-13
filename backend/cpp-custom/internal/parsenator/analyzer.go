package parsenator

import (
	"cpp-custom/internal/lexinator"
	"cpp-custom/internal/semanthoid"
	"cpp-custom/logger"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

// analyzer structure
type Analyzer struct {
	scanner    lexinator.Scanner
	isPrepared bool
	// the output stream of error messages
	writer io.Writer
	// interpretation flag
	IFlag int
}

// preparing the analyzer
func Preparing(srsFileName string, scannerErrWriter io.Writer, analyzerErrWriter io.Writer) (Analyzer, error) {
	// file loggers preparing
	if os.Getenv("ENABLE_LOGGING") == "true" {
		loggers := make(map[string]string)
		loggers["memory_l"] = "memory"
		loggers["procedures_tree_l"] = "procedures_tree"
		loggers["tree_l"] = "tree"
		loggers["stack_l"] = "stack"
		loggers["debug_l"] = "debug"
		err := logger.Init(loggers)
		if err != nil {
			panic("error logger initializing")
		}
	}
	// syntax analyzer
	A := Analyzer{writer: analyzerErrWriter, IFlag: 0x0000}
	// lexical analyzer
	scanner, err := lexinator.ScannerInitializing(srsFileName, scannerErrWriter)
	if err != nil {
		return A, err
	}
	A.scanner = scanner
	// symantic tree
	if true {
		// TODO: add saving parse tree and proper node preparing
		semanthoid.Root = nil
		semanthoid.Current = nil
		semanthoid.ProcRoot = nil
		semanthoid.CurrentProc = nil
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

// interpretation flag setting
func (A *Analyzer) resetIFlag() {
	A.IFlag = 0x000
}

func (A *Analyzer) setGlobalDescriptionMode() {
	A.IFlag = 0x1000
}

func (A *Analyzer) setProcedureExecutionMode() {
	A.IFlag = 0x0100
}

func (A *Analyzer) setProcedureDescriptionMode() {
	A.IFlag = 0x0010
}

func (A *Analyzer) isInterpretingGlobalDescription() bool {
	return A.IFlag&0x1000 != 0
}

func (A *Analyzer) isInterpretingProcedureExecution() bool {
	return A.IFlag&0x0100 != 0
}

func (A *Analyzer) isInterpretingExpression() bool {
	return A.IFlag&0x1100 != 0
}

func (A *Analyzer) isInterpretingProcedureDescription() bool {
	return A.IFlag&0x0010 != 0
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
			A.setProcedureDescriptionMode()
			procedureIdentifier := A.procedureDescription()
			if procedureIdentifier == "main" { // if run program execution
				mainProc, err := semanthoid.LoadProcedure("main", []int{}) // loading main procedure
				if err != nil {
					A.printPanicError(err.Error())
				}
				A.setProcedureExecutionMode()
				A.scanner.RestorePosValues(mainProc.ProcTextPos, mainProc.ProcLine, mainProc.ProcLinePos)
				A.procedureDescription() // run main execution
				break                    // stop parsing
			}
		} else if lexType == lexinator.Long ||
			lexType == lexinator.Short ||
			lexType == lexinator.Int ||
			lexType == lexinator.Bool ||
			lexType == lexinator.Const { // <описание>
			A.scanner.RestorePosValues(textPos, line, linePos)
			A.setGlobalDescriptionMode()
			A.description()
		} else if lexType != lexinator.Semicolon { // then must be ';'
			A.printPanicError("invalid lexeme '" + lex + "', expected ';'")
		}
		A.resetIFlag()
		textPos, line, linePos = A.scanner.StorePosValues()
		lexType, lex = A.scanner.Scan()
	}
	A.printError("there are no syntax level errors")
	logger.Log("tree_l", "endpoint tree:\n"+semanthoid.TreeToString())
	return nil
}

// <описание параметров> ->
// return paramsCount, paramsTypes and paramsIdentifiers
func (A *Analyzer) parameterDescription() (int, []int, []string) {
	var textPos, line, linePos int
	var lexType int
	var lex string
	isFirst, paramsCount := true, 0
	var paramsTypes []int
	var paramsIdentifiers []string
	for isFirst || lexType == lexinator.Comma {
		isFirst = false
		paramsTypes = append(paramsTypes, A._type())
		paramsCount++
		lexType, lex = A.scanner.Scan()
		if lexType != lexinator.Id {
			A.printPanicError("'" + lex + "' is not an identifier")
		}
		paramsIdentifiers = append(paramsIdentifiers, lex)
		textPos, line, linePos = A.scanner.StorePosValues()
		lexType, lex = A.scanner.Scan()
	}
	A.scanner.RestorePosValues(textPos, line, linePos)
	return paramsCount, paramsTypes, paramsIdentifiers
}

// <параметры> -> идентификатор | константа U , <параметры> | e
// return params count
func (A *Analyzer) parameters() []int {
	var paramsValues []int
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
		if A.isInterpretingProcedureExecution() {
			var value int
			var err error
			if lexType == lexinator.IntConst {
				value, err = strconv.Atoi(lex)
				if err != nil {
					A.printPanicError("bad integer constant '" + lex + "'")
				}
			} else { // identifier
				dataNode := semanthoid.FindDataUpFromCurrent(lex)
				if dataNode == nil {
					A.printPanicError("'" + lex + "' undefined")
				} else {
					value = dataNode.DataValue.DataAsInt
				}
			}
			paramsValues = append(paramsValues, value)
		}
		textPos, line, linePos = A.scanner.StorePosValues()
		lexType, lex = A.scanner.Scan()
	}
	A.scanner.RestorePosValues(textPos, line, linePos)
	return paramsValues
}

// <эл.выр.> -> (<выражение>) | идентификатор | константа
// returns dataTypeLabel, dataAsInt value
func (A *Analyzer) simplestExpr() int {
	value := 0

	lexType, lex := A.scanner.Scan()
	if lexType == lexinator.OpeningBracket { // ( <выражение> )
		value = A.expression()
		lexType, lex = A.scanner.Scan()
		if lexType != lexinator.ClosingBracket {
			A.printPanicError("invalid lexeme '" + lex + "', expected ')'")
		}
		return value
	} else if lexType != lexinator.Id && lexType != lexinator.IntConst {
		A.printPanicError("'" + lex + "' not allowed in the expression")
	} else if A.isInterpretingExpression() {
		if lexType == lexinator.Id {
			dataNode := semanthoid.FindDataUpFromCurrent(lex)
			if dataNode == nil {
				A.printPanicError("'" + lex + "' undefined")
			} else {
				return dataNode.DataValue.DataAsInt
			}
		} else {
			intConstVal, _ := strconv.Atoi(lex)
			return intConstVal
		}
	}
	return value
}

// <множитель> -> <эл.выр.> U e | * U / U % <эл.выр.>
//
func (A *Analyzer) multiplier() int {
	value := A.simplestExpr() // <эл.выр.>
	textPos, line, linePos := A.scanner.StorePosValues()
	lexType, _ := A.scanner.Scan()
	for lexType == lexinator.Mul || lexType == lexinator.Div || lexType == lexinator.Mod {
		value2 := A.simplestExpr()
		switch lexType {
		case lexinator.Mul:
			value *= value2
			break
		case lexinator.Div:
			value /= value2
			break
		case lexinator.Mod:
			value %= value2
			break
		}
		textPos, line, linePos = A.scanner.StorePosValues()
		lexType, _ = A.scanner.Scan()
	}
	A.scanner.RestorePosValues(textPos, line, linePos)
	return value
}

// <процедура> -> идентификатор ( ) | идентификатор ( <параметры> )
// return procedure id, paramsCount and paramsValues
func (A *Analyzer) procedure() (string, []int) {
	var procIdentifier string
	var paramsValues []int

	lexType, lex := A.scanner.Scan()
	if lexType != lexinator.Id {
		A.printPanicError("'" + lex + "' is not an identifier")
	}
	procIdentifier = lex
	lexType, lex = A.scanner.Scan()
	if lexType != lexinator.OpeningBracket {
		A.printPanicError("invalid lexeme '" + lex + "', expected '('")
	}
	textPos, line, linePos := A.scanner.StorePosValues()
	lexType, lex = A.scanner.Scan()
	if lexType != lexinator.ClosingBracket {
		A.scanner.RestorePosValues(textPos, line, linePos)
		paramsValues = A.parameters()
		lexType, lex = A.scanner.Scan()
		if lexType != lexinator.ClosingBracket {
			A.printPanicError("invalid lexeme '" + lex + "', expected ')'")
		}
	}
	return procIdentifier, paramsValues
}

// <слагаемое> -> <множитель> U +- | e
// returns DatatypeLabel and DataTypeValue
func (A *Analyzer) summand() int {
	value := A.multiplier() // <множитель>
	textPos, line, linePos := A.scanner.StorePosValues()
	lexType, _ := A.scanner.Scan()
	for lexType == lexinator.Plus || lexType == lexinator.Minus {
		value2 := A.multiplier()
		switch lexType {
		case lexinator.Plus:
			value += value2
			break
		case lexinator.Minus:
			value -= value2
			break
		}
		textPos, line, linePos = A.scanner.StorePosValues()
		lexType, _ = A.scanner.Scan()
	}
	A.scanner.RestorePosValues(textPos, line, linePos)
	return value
}

// <выражение> -> + | - | e U <слагаемое> + +- == <= >= < > <слагаемое> | e
func (A *Analyzer) expression() int {
	textPos, line, linePos := A.scanner.StorePosValues()
	lexType, _ := A.scanner.Scan()
	sign := 1
	if lexType != lexinator.Plus && lexType != lexinator.Minus {
		A.scanner.RestorePosValues(textPos, line, linePos)
	} else if lexType == lexinator.Minus {
		sign = -1
	}
	value := A.summand() * sign // <слагаемое> * sign
	textPos, line, linePos = A.scanner.StorePosValues()
	lexType, _ = A.scanner.Scan()
	for lexType == lexinator.Plus || lexType == lexinator.Minus ||
		lexType == lexinator.Equ || lexType == lexinator.LessEqu || lexType == lexinator.MoreEqu ||
		lexType == lexinator.Less || lexType == lexinator.More {
		value2 := A.multiplier()
		switch lexType {
		case lexinator.Plus:
			value += value2
			break
		case lexinator.Minus:
			value -= value2
			break
		case lexinator.Equ:
			value = semanthoid.GoBoolToInt(value == value2)
			break
		case lexinator.LessEqu:
			value = semanthoid.GoBoolToInt(value <= value2)
			break
		case lexinator.MoreEqu:
			value = semanthoid.GoBoolToInt(value >= value2)
			break
		case lexinator.Less:
			value = semanthoid.GoBoolToInt(value < value2)
			break
		case lexinator.More:
			value = semanthoid.GoBoolToInt(value > value2)
			break
		}
		textPos, line, linePos = A.scanner.StorePosValues()
		lexType, _ = A.scanner.Scan()
	}
	A.scanner.RestorePosValues(textPos, line, linePos)
	return value
}

// <оператор for>
func (A *Analyzer) forOperator() (int, string, *semanthoid.DataTypeValue) {
	counterType := semanthoid.Error
	counterIdentifier := ""
	textPos, line, linePos := A.scanner.StorePosValues()
	lexType, lex := A.scanner.Scan()
	if lexType != lexinator.Semicolon { // if not empty
		if lexType != lexinator.Id { // if type
			A.scanner.RestorePosValues(textPos, line, linePos)
			counterType = A._type()
			lexType, lex = A.scanner.Scan()
			if lexType != lexinator.Id {
				A.printPanicError("'" + lex + "' is not an identifier")
			}
			counterIdentifier = lex
		}
		lexType, lex = A.scanner.Scan()
		if lexType != lexinator.Assignment {
			A.printPanicError("invalid lexeme '" + lex + "', expected '='")
		}
		A.expression()
	}
	return counterType, counterIdentifier, semanthoid.GetDefaultDataValue() // TODO: add counter value calculation
}

// <присваивание> -> идентификатор = <выражение>
func (A *Analyzer) assigment() {
	lexType, lex := A.scanner.Scan()
	if lexType != lexinator.Id {
		A.printPanicError("'" + lex + "' is not an identifier")
	}
	identifier := lex
	lexType, lex = A.scanner.Scan()
	if lexType != lexinator.Assignment {
		A.printPanicError("invalid lexeme '" + lex + "', expected '='")
	}
	value := A.expression()
	if A.isInterpretingProcedureExecution() {
		variable := semanthoid.FindVariableUpFromCurrent(identifier)
		if variable == nil {
			A.printPanicError("undefined variable '" + identifier + "'")
		}
		variable.DataValue = semanthoid.ConvertToDataTypeValue(variable.DataTypeLabel, value)
	}
}

// <переменная> -> идентификатор U e | = <выражение>
func (A *Analyzer) variable() (string, int) {
	lexType, lex := A.scanner.Scan()
	if lexType != lexinator.Id {
		A.printPanicError("'" + lex + "' is not an identifier")
	}
	identifier := lex // variable identifier
	value := 0
	textPos, line, linePos := A.scanner.StorePosValues()
	lexType, lex = A.scanner.Scan()
	if lexType == lexinator.Assignment {
		value = A.expression()
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
// receives current right subtree root
func (A *Analyzer) constants() {
	lexType, lex := A.scanner.Scan()
	if lexType != lexinator.Const {
		A.printPanicError("invalid lexeme '" + lex + "', expected 'const'")
	}
	constsType := A._type()
	var textPos, line, linePos int
	isFirst := true
	for isFirst || lexType == lexinator.Comma {
		// parsing
		isFirst = false
		lexType, lex := A.scanner.Scan()
		if lexType != lexinator.Id {
			A.printPanicError("'" + lex + "' is not an identifier")
		}
		identifier := lex
		lexType, lex = A.scanner.Scan()
		if lexType != lexinator.Assignment {
			A.printPanicError("invalid lexeme '" + lex + "', expected '='")
		}
		value := A.expression()
		// interpreting
		// value preparation
		dataValue := semanthoid.ConvertToDataTypeValue(constsType, value)
		if A.isInterpretingGlobalDescription() {
			err := semanthoid.CreateGlobalDescription(semanthoid.Constant, identifier, constsType, dataValue)
			if err != nil {
				A.printPanicError(err.Error())
			}
		} else if A.isInterpretingProcedureExecution() {
			err := semanthoid.CreateLocalDescription(semanthoid.Constant, identifier, constsType, dataValue)
			if err != nil {
				A.printPanicError(err.Error())
			}
		}
		textPos, line, linePos = A.scanner.StorePosValues()
		lexType, lex = A.scanner.Scan()
	}
	A.scanner.RestorePosValues(textPos, line, linePos)
}

// <переменные> -> long int | short int | int | bool U e | <присваивание> U e | , <присваивание>
// receives current right subtree root
// returns variables subtree
func (A *Analyzer) variables() {
	varsType := A._type()
	var lexType, textPos, line, linePos int
	isFirst := true
	for isFirst || lexType == lexinator.Comma {
		isFirst = false
		identifier, value := A.variable()
		// value preparation
		dataValue := semanthoid.ConvertToDataTypeValue(varsType, value)
		if A.isInterpretingGlobalDescription() {
			err := semanthoid.CreateGlobalDescription(semanthoid.Variable, identifier, varsType, dataValue)
			if err != nil {
				A.printPanicError(err.Error())
			}
		} else if A.isInterpretingProcedureExecution() {
			err := semanthoid.CreateLocalDescription(semanthoid.Variable, identifier, varsType, dataValue)
			if err != nil {
				A.printPanicError(err.Error())
			}
		}
		textPos, line, linePos = A.scanner.StorePosValues()
		lexType, _ = A.scanner.Scan()
	}
	A.scanner.RestorePosValues(textPos, line, linePos)
}

// <оператор> -> <составной оператор> | <for> | <процедура> ; | <присваивание>; | ;
// returns isReturn
func (A *Analyzer) operator() bool {
	isReturn := false
	textPos, line, linePos := A.scanner.StorePosValues()
	lexType, lex := A.scanner.Scan()
	if lexType == lexinator.OpeningBrace { // составной оператор
		A.scanner.RestorePosValues(textPos, line, linePos)
		A.compositeOperator(true) // with fork creating
	} else if lexType == lexinator.For { // for
		A.scanner.RestorePosValues(textPos, line, linePos)
		A._for()
	} else if lexType == lexinator.Id || lexType == lexinator.Return { // процедура или присваивание
		if lexType == lexinator.Return && A.isInterpretingProcedureExecution() {
			return true // terminate procedure and return terminate flag
		} else if lexType == lexinator.Id {
			lexType, lex = A.scanner.Scan()
			if lexType == lexinator.OpeningBracket { // процедура
				A.scanner.RestorePosValues(textPos, line, linePos)
				procIdentifier, paramsValues := A.procedure()
				if A.isInterpretingProcedureExecution() { // need run proc
					// steps:
					// 0. save scanner position
					// 1. save current procedure context
					// 2. load procedure for execution
					// 3. set scanner position for procedure
					// 4. run procedureDescription for procedure execution (flag already set)
					// 5. load procedure context back
					// 6. set old scanner position
					logger.Log("debug_l", "\n"+semanthoid.TreeToString())
					logger.Log("debug_l", semanthoid.Current.ToString())
					textPos, line, linePos = A.scanner.StorePosValues() // 0
					err := semanthoid.SaveProcedureContext()            // 1
					if err != nil {
						A.printPanicError(err.Error())
					}
					logger.Log("debug_l", "\n"+semanthoid.TreeToString())
					logger.Log("debug_l", semanthoid.Current.ToString())
					proc, err := semanthoid.LoadProcedure(procIdentifier, paramsValues) // 2
					if err != nil {
						A.printPanicError(err.Error())
					}
					logger.Log("debug_l", "\n"+semanthoid.TreeToString())
					logger.Log("debug_l", semanthoid.Current.ToString())
					A.scanner.RestorePosValues(proc.ProcTextPos, proc.ProcLine, proc.ProcLinePos) // 3
					A.procedureDescription()                                                      // 4
					logger.Log("debug_l", semanthoid.Current.ToString())
					err = semanthoid.LoadProcedureContext() // 5
					if err != nil {
						A.printPanicError(err.Error())
					}
					logger.Log("debug_l", "\n"+semanthoid.TreeToString())
					logger.Log("debug_l", semanthoid.Current.ToString())
					A.scanner.RestorePosValues(textPos, line, linePos) // 6
				}
			} else if lexType == lexinator.Assignment { // присваивание
				A.scanner.RestorePosValues(textPos, line, linePos)
				A.assigment()
			} else {
				A.printPanicError("'" + lex + "' is not an procedure or assigment")
			}
		}
		lexType, lex = A.scanner.Scan()
		if lexType != lexinator.Semicolon {
			A.printPanicError("invalid lexeme '" + lex + "', expected ';'")
		}
	} else if lexType != lexinator.Semicolon {
		A.printPanicError("invalid lexeme '" + lex + "', expected ';'")
	}
	return isReturn
}

// <описание процедуры> -> void идентификатор ( U <описание параметров> | e U ) <составной оператор>
// return proc identifier
func (A *Analyzer) procedureDescription() string {
	// procedure
	// procedure node data
	var procedureIdentifier string
	var paramsCount int
	var paramsTypes []int
	var paramsIdentifiers []string
	procTextPos, procLine, procLinePos := A.scanner.StorePosValues()
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
		paramsCount, paramsTypes, paramsIdentifiers = A.parameterDescription()
		if procedureIdentifier == "main" && paramsCount > 0 {
			A.printPanicError("'main' must not have input arguments")
		}
		lexType, lex = A.scanner.Scan()
		if lexType != lexinator.ClosingBracket {
			A.printPanicError("invalid lexeme '" + lex + "', expected ')'")
		}
	} else if lexType != lexinator.ClosingBracket { // )
		A.printPanicError("invalid lexeme '" + lex + "', expected ')'")
	}
	if A.isInterpretingProcedureDescription() {
		err := semanthoid.AddProcedureDescription(procedureIdentifier, paramsCount, paramsTypes, paramsIdentifiers,
			procTextPos, procLine, procLinePos)
		if err != nil {
			A.printPanicError(err.Error())
		}
	}
	A.compositeOperator(false) // without fork creating
	return procedureIdentifier
}

// <составной оператор> -> { <операторы и описания> }
func (A *Analyzer) compositeOperator(isNeedCreateFork bool) {
	if A.isInterpretingProcedureExecution() && isNeedCreateFork {
		semanthoid.CreateFork()
	}

	lexType, lex := A.scanner.Scan()
	if lexType != lexinator.OpeningBrace {
		A.printPanicError("invalid lexeme '" + lex + "', expected '{'")
	}
	A.operatorsAndDescriptions() // операторы и описания
	lexType, lex = A.scanner.Scan()
	if lexType != lexinator.ClosingBrace && !A.isInterpretingProcedureExecution() { // next parsing if procedure not execute now
		A.printPanicError("invalid lexeme '" + lex + "', expected '}'")
	}

	if A.isInterpretingProcedureExecution() {
		err := semanthoid.ClearCurrentRightSubTree()
		if err != nil {
			A.printPanicError(err.Error())
		}
	}
}

// <операторы и описания> -> e | <операторы> U e | <операторы и описания>  | <описания> + e | <операторы и описания>
func (A *Analyzer) operatorsAndDescriptions() {
	for {
		textPos, line, linePos := A.scanner.StorePosValues()
		lexType, _ := A.scanner.Scan()
		A.scanner.RestorePosValues(textPos, line, linePos)
		if lexType == lexinator.OpeningBrace || lexType == lexinator.For || lexType == lexinator.Return ||
			lexType == lexinator.Id || lexType == lexinator.Semicolon { // if operator
			if A.operator() && A.isInterpretingProcedureExecution() {
				return // terminate execution
			}
		} else if lexType == lexinator.Long || lexType == lexinator.Short ||
			lexType == lexinator.Int || lexType == lexinator.Bool || lexType == lexinator.Const { // if description
			A.description()
		} else { // e
			break
		}
	}
}

// <описание> -> <переменные>; | <константы>;
// receives current nearest description subtree
// returns description subtree node
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
		return semanthoid.BoolType
	} else {
		return semanthoid.IntType
	}
}
