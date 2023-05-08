package ll1

import (
	"cpp-custom/internal/il"
	"cpp-custom/internal/lexinator"
	"cpp-custom/internal/stree"
	"cpp-custom/logger"
	"errors"
	"fmt"
	"github.com/golang-collections/collections/stack"
	"io"
	"strings"
)

const (
	LlTableExcelPath      = "./internal/ll1/specifications/LL_table.xlsx"
	LlTableExcelSheet     = "Sheet1"
	NonTerminalsCellCount = 36
	TerminalsCellCount    = 29
)

const (
	StackL      = "stack_l"
	RuleL       = "rule_l"
	TreeL       = "tree_l"
	OperationsL = "operations_l"
)

// LlChecker struct with methods
type LlChecker struct {
	scanner lexinator.Scanner // use scanner
	llTable *LlTable          // llTable for checking
	stack   *stack.Stack
	//
	root       *stree.Root      // root of sym tree
	operations *il.Intermediate // operations of triad
	// the output stream of error messages
	writer io.Writer
}

func CreateLlChecker(srcFileName string, scannerErrWriter io.Writer, llCheckerErrWriter io.Writer) (*LlChecker, error) {
	// file loggers preparing
	loggers := make(map[string]string)
	loggers[StackL] = "stack"
	loggers[RuleL] = "rule"
	loggers[TreeL] = "tree"
	loggers[OperationsL] = "operations"
	err := logger.Init(loggers)
	if err != nil {
		panic("error logger initializing")
	}

	// read llTable
	llTable, err := ReadLLTableFromExcel(LlTableExcelPath, LlTableExcelSheet, NonTerminalsCellCount, TerminalsCellCount)
	if err != nil {
		return nil, err
	}

	// create scanner
	scanner, err := lexinator.ScannerInitializing(srcFileName, scannerErrWriter)
	if err != nil {
		return nil, err
	}

	// create checker
	checker := new(LlChecker)
	checker.scanner = scanner
	checker.llTable = llTable
	checker.root = stree.NewRoot()
	checker.operations = il.NewIntermediateDomain(1000) // temporary it will be 1000 operations
	checker.writer = llCheckerErrWriter
	checker.stack = stack.New()
	return checker, nil
}

// the results of the error message
func (c *LlChecker) printPanicError(msg string) {
	textPos, line, linePos := c.scanner.StorePosValues()
	_, err := fmt.Fprintf(c.writer, "error: %s position: %d line: %d line position: %d\n", msg, textPos, line, linePos)
	if err != nil {
		panic(err)
	}
	// maybe temporarily
	panic(errors.New("completed with an error. see the error logs")) // critical completion
}

func (c *LlChecker) printError(msg string) {
	textPos, line, linePos := c.scanner.StorePosValues()
	_, err := fmt.Fprintf(c.writer, "error: %s position: %d line: %d line position: %d\n", msg, textPos, line, linePos)
	if err != nil {
		panic(err)
	}
}

func (c *LlChecker) pushToStack(value string) {
	c.stack.Push(value)
	logger.Log(StackL, "pushed: "+value)
}

func (c *LlChecker) extractFromStack() string {
	logger.Log(StackL, "extracted: "+c.watchStack())
	return fmt.Sprintf("%v", c.stack.Pop())
}

func (c *LlChecker) watchStack() string {
	return fmt.Sprintf("%v", c.stack.Peek())
}

func (c *LlChecker) isTerminalOnStackTop() bool {
	return c.llTable.IsTerminal(c.watchStack())
}

func (c *LlChecker) isNonTerminalOnStackTop() bool {
	return c.llTable.IsNonTerminal(c.watchStack())
}

func (c *LlChecker) isOperationalOnStackTop() bool {
	return c.llTable.IsOperational(c.watchStack())
}

func (c *LlChecker) stackToString() string {
	builder := strings.Builder{}
	tmpStack := stack.New()
	for c.stack.Len() > 0 {
		val := c.extractFromStack()
		builder.WriteString(val + "\n")
		tmpStack.Push(val)
	}
	for tmpStack.Len() > 0 {
		c.stack.Push(tmpStack.Pop())
	}
	return builder.String()
}

func (c *LlChecker) TreeToString() string {
	return c.root.AsString()
}

func (c *LlChecker) IntermediateCode() string {
	return c.operations.IntermediateAsString()
}

func (c *LlChecker) MakeLkAnalyze() {
	// 1. add first non terminal to stack
	c.stack.Push(c.llTable.GetFirstNonTerminal())

	// 1.1. prepare context and decorator
	ctx := &context{"", "", "", ""}
	// this need to scan identifiers and constants literals
	scanDecorator := func() (int, string, string) {
		lexType, lex := c.scanner.Scan()
		syntaxLex := lex
		if lexType == lexinator.Id || lexType == lexinator.Main {
			syntaxLex = "IDENTITY"
			ctx.saveIdentity(lex)
		} else if lexType == lexinator.IntConst {
			syntaxLex = "CONSTANT"
			ctx.saveConstant(lex)
		} else if lexType == lexinator.Void ||
			lexType == lexinator.Short ||
			lexType == lexinator.Long ||
			lexType == lexinator.Int ||
			lexType == lexinator.Bool {
			ctx.saveType(lex)
		}
		return lexType, lex, syntaxLex
	}

	// 2. parse input
	lexType, lex, sLex := scanDecorator()
	for lexType != lexinator.End || c.isOperationalOnStackTop() {
		if c.isTerminalOnStackTop() { // terminal on top
			if c.watchStack() == sLex {
				// same terminals on top and on input
				c.extractFromStack()
				lexType, lex, sLex = scanDecorator()
			} else {
				c.printPanicError("expected symbol '" + c.watchStack() + "'")
			}
		} else if c.isNonTerminalOnStackTop() { // not terminal on top
			// find terminal
			nonTerminalOnTop := c.watchStack()
			reversedRule, isFound := c.llTable.Table[nonTerminalOnTop][sLex]
			if isFound && reversedRule != "EPSILON" {
				// rule is found in ll table then apply this rule
				logger.Log(RuleL, "apply reversed rule "+reversedRule+" from "+nonTerminalOnTop)
				c.extractFromStack() // pop non terminal
				ruleParts := strings.Split(reversedRule, " ")
				for _, rulePart := range ruleParts {
					c.pushToStack(rulePart)
				}
			} else if reversedRule == "EPSILON" || c.llTable.IsEpsilonRuleExists(nonTerminalOnTop) {
				logger.Log(RuleL, "exists EPSILON rule for "+nonTerminalOnTop+" and input "+lex+" (syntax lex "+sLex+")")
				c.extractFromStack() // pop non terminal
			} else {
				expected := c.llTable.GetTerminalsListByNonTerminal(nonTerminalOnTop)
				c.printPanicError("Bad input '" + lex + "' in '" + nonTerminalOnTop + "', on of " + strings.Join(expected, " ") + " expected")
			}
		} else if c.isOperationalOnStackTop() { // operational symbol on stack top
			operationalSymbol := c.extractFromStack()
			logger.Info.Println("Operational symbol on stack top: " + operationalSymbol)
			if err := runOperational(operationalSymbol, c.root, c.operations, ctx); err != nil {
				c.printPanicError(err.Error())
			}
		} else {
			c.printPanicError("Bad value on stack top: " + c.watchStack())
		}
	}

	logger.Log(StackL, "stack state at the end:\n"+c.stackToString())
	logger.Log(TreeL, "tree:\n"+c.root.AsString())
	logger.Log(OperationsL, "operations in intermediate:\n"+c.operations.IntermediateAsString())
	c.printError("there are no ll-level errors")
}
