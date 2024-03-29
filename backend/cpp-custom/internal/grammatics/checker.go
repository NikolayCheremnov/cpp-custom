package main

import (
	"cpp-custom/internal/lexinator"
	"cpp-custom/logger"
	"errors"
	"fmt"
	"github.com/golang-collections/collections/stack"
	"io"
	"strings"
)

const (
	LlTableExcelPath      = "./specifications/LL_table.xlsx"
	LlTableExcelSheet     = "Sheet1"
	NonTerminalsCellCount = 36
	TerminalsCellCount    = 29
	LogDir                = "./tdata/logs/"
)

const (
	StackL = "stack_l"
	RuleL  = "rule_l"
)

// LlChecker struct with methods
type LlChecker struct {
	scanner lexinator.Scanner // use scanner
	llTable *LlTable          // llTable for checking
	stack   *stack.Stack
	// the output stream of error messages
	writer io.Writer
}

func CreateLlChecker(srcFileName string, scannerErrWriter io.Writer, llCheckerErrWriter io.Writer) (*LlChecker, error) {
	// file loggers preparing
	loggers := make(map[string]string)
	loggers[StackL] = "stack"
	loggers[RuleL] = "rule"
	err := logger.InitWithCustomLogDir(loggers, LogDir)
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

func (c *LlChecker) isTerminalOnStackTop() bool {
	return c.llTable.IsTerminal(c.watchStack())
}

func (c *LlChecker) extractFromStack() string {
	logger.Log(StackL, "extracted: "+c.watchStack())
	return fmt.Sprintf("%v", c.stack.Pop())
}

func (c *LlChecker) watchStack() string {
	return fmt.Sprintf("%v", c.stack.Peek())
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

func (c *LlChecker) MakeLkAnalyze() {
	// 1. add first non terminal to stack
	c.stack.Push(c.llTable.GetFirstNonTerminal())

	// 2. parse input

	// this need to scan identifiers and constants literals
	scanDecorator := func() (int, string, string) {
		lexType, lex := c.scanner.Scan()
		syntaxLex := lex
		if lexType == lexinator.Id || lexType == lexinator.Main {
			syntaxLex = "IDENTITY"
		} else if lexType == lexinator.IntConst {
			syntaxLex = "CONSTANT"
		}
		return lexType, lex, syntaxLex
	}

	lexType, lex, sLex := scanDecorator()
	for lexType != lexinator.End {
		if c.isTerminalOnStackTop() { // terminal on top
			if c.watchStack() == sLex {
				// same terminals on top and on input
				c.extractFromStack()
				lexType, lex, sLex = scanDecorator()
			} else {
				c.printPanicError("expected symbol '" + c.watchStack() + "'")
			}
		} else { // not terminal on top
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
		}
	}
	c.printError("there are no ll-level errors")
}
