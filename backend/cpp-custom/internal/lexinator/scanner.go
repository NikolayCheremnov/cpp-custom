package lexinator

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

// scanner structure
type Scanner struct {
	sourceModule string // source code
	textPos      int    // position in text
	line         int    // line number
	linePos      int    // position in line

	// the output stream of error messages
	writer io.Writer
}

// initialize function
func ScannerInitializing(srsFileName string, scannerErrWriter io.Writer) (S Scanner, err error) {
	S = Scanner{sourceModule: "undefined", textPos: 0, line: 0, linePos: 0, writer: scannerErrWriter}
	err = S.GetData(srsFileName)
	if err != nil {
		return S, err
	}
	return S, nil
}

// store and restore functions
func (S *Scanner) StorePosValues() (int, int, int) {
	return S.textPos, S.line, S.linePos
}

func (S *Scanner) RestorePosValues(textPos int, line int, linePos int) {
	S.textPos = textPos
	S.line = line
	S.linePos = line
}

// getters and setters
func (S *Scanner) GetTextPos() int {
	return S.textPos
}

func (S *Scanner) GetLine() int {
	return S.line
}

func (S *Scanner) GetLinePos() int {
	return S.linePos
}

func (S *Scanner) SetTextPos(val int) {
	S.textPos = val
}

func (S *Scanner) SetLine(val int) {
	S.line = val
}

func (S *Scanner) SetLinePos(val int) {
	S.linePos = val
}

// read source module
func (S *Scanner) GetData(fname string) (err error) {

	bytes, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}
	S.sourceModule = string(bytes)
	S.sourceModule = strings.ReplaceAll(S.sourceModule, "\r", "") // remove excess /r
	if len(S.sourceModule) > MaxModuleLen {
		return errors.New("too long module")
	}
	S.sourceModule += "\000" // terminate byte
	return nil
}

// the results of the error message
func (S *Scanner) printError(msg string) {
	_, err := fmt.Fprintf(S.writer, "error: %s postion: %d line: %d line position: %d\n", msg, S.textPos, S.line, S.linePos)
	if err != nil {
		panic(err)
	}
}

// scanner function
// input parameter: source module text, position in text, line index, position in line,
// output: lexType, lexImage, error and changed parameters
func (S *Scanner) Scan() (lexType int, lex string) {
	incPos := func() { S.textPos++; S.linePos++ }
	incLine := func() { S.line++; S.linePos = 0 }

	for {
		// preparing
		lex = "" // lexeme image

		// ignored symbols passing
		for S.sourceModule[S.textPos] == ' ' || S.sourceModule[S.textPos] == '\n' || S.sourceModule[S.textPos] == '\t' {
			if S.sourceModule[S.textPos] == '\n' {
				incLine()
			}
			incPos()
		}

		// next conditions
		// 1) is a latin letter
		if (S.sourceModule[S.textPos] >= 'a' && S.sourceModule[S.textPos] <= 'z') ||
			(S.sourceModule[S.textPos] >= 'A' && S.sourceModule[S.textPos] <= 'Z') {
			for (S.sourceModule[S.textPos] >= 'a' && S.sourceModule[S.textPos] <= 'z') ||
				(S.sourceModule[S.textPos] >= 'A' && S.sourceModule[S.textPos] <= 'Z') ||
				(S.sourceModule[S.textPos] >= '0' && S.sourceModule[S.textPos] <= '9') {
				lex += S.sourceModule[S.textPos : S.textPos+1] // add symbol to lexeme
				incPos()
				if len(lex) > MaxLexLen {
					S.printError("too long lexeme")
					return Err, lex
				}
			}
			// id was retrieved => check keyword
			switch lex {
			case "main":
				return Main, lex
			case "int":
				return Int, lex
			case "short":
				return Short, lex
			case "long":
				return Long, lex
			case "bool":
				return Bool, lex
			case "for":
				return For, lex
			case "const":
				return Const, lex
			case "void":
				return Void, lex
			default:
				return Id, lex
			}
		}

		// 2) is a digital
		if S.sourceModule[S.textPos] >= '0' && S.sourceModule[S.textPos] <= '9' {
			for S.sourceModule[S.textPos] >= '0' && S.sourceModule[S.textPos] <= '9' {
				lex += S.sourceModule[S.textPos : S.textPos+1]
				incPos()
				if len(lex) > MaxLexLen {
					S.printError("too long lexeme")
					return Err, lex
				}
			}
			return IntConst, lex
		}

		// 3) single special symbols
		if S.sourceModule[S.textPos] == ';' {
			incPos()
			return Semicolon, ";"
		} else if S.sourceModule[S.textPos] == ',' {
			incPos()
			return Comma, ","
		} else if S.sourceModule[S.textPos] == '(' {
			incPos()
			return OpeningBracket, "("
		} else if S.sourceModule[S.textPos] == ')' {
			incPos()
			return ClosingBracket, ")"
		} else if S.sourceModule[S.textPos] == '{' {
			incPos()
			return OpeningBrace, "{"
		} else if S.sourceModule[S.textPos] == '}' {
			incPos()
			return ClosingBrace, "}"
		} else if S.sourceModule[S.textPos] == '+' {
			incPos()
			return Plus, "+"
		} else if S.sourceModule[S.textPos] == '-' {
			incPos()
			return Minus, "-"
		} else if S.sourceModule[S.textPos] == '*' {
			incPos()
			return Mul, "*"
		} else if S.sourceModule[S.textPos] == '%' {
			incPos()
			return Mod, "%"
		}

		// 4) /
		if S.sourceModule[S.textPos] == '/' {
			lex = S.sourceModule[S.textPos : S.textPos+1]
			incPos()
			if S.sourceModule[S.textPos] == '/' { // line comments
				for S.sourceModule[S.textPos] != '\n' {
					if S.sourceModule[S.textPos] == '\000' {
						return End, "null_char"
					}
					incPos()
				}
				incPos()
				incLine()
				continue // go to start again
			} else if S.sourceModule[S.textPos] == '*' { // multiline comment
				incPos()
				for S.sourceModule[S.textPos] != '*' || (S.sourceModule[S.textPos] == '*' && S.sourceModule[S.textPos+1] != '/') {
					if S.sourceModule[S.textPos] == '\n' {
						incPos()
						incLine()
					} else if S.sourceModule[S.textPos] == '\000' {
						S.printError("unclosed comment in end")
						return Err, "/*"
					} else {
						incPos()
					}
				}
				incPos() // eat *
				incPos() // eat /
				continue
			} else {
				return Div, lex
			}
		}

		// <
		if S.sourceModule[S.textPos] == '<' {
			lex = S.sourceModule[S.textPos : S.textPos+1]
			incPos()
			if S.sourceModule[S.textPos] == '=' {
				lex += S.sourceModule[S.textPos : S.textPos+1]
				incPos()
				return LessEqu, lex
			} else {
				return Less, lex
			}
		}

		// >
		if S.sourceModule[S.textPos] == '>' {
			lex = S.sourceModule[S.textPos : S.textPos+1]
			incPos()
			if S.sourceModule[S.textPos] == '=' {
				lex += S.sourceModule[S.textPos : S.textPos+1]
				incPos()
				return MoreEqu, lex
			} else {
				return More, lex
			}
		}

		// =
		if S.sourceModule[S.textPos] == '=' {
			lex = S.sourceModule[S.textPos : S.textPos+1]
			incPos()
			if S.sourceModule[S.textPos] == '=' {
				lex += S.sourceModule[S.textPos : S.textPos+1]
				incPos()
				return Equ, lex
			} else {
				return Assignment, lex
			}
		}

		// module end
		if S.sourceModule[S.textPos] == '\000' {
			return End, "null_char"
		}

		// error on input
		incPos()
		S.printError("invalid input symbol: " + S.sourceModule[S.textPos-1:S.textPos])
		return Err, S.sourceModule[S.textPos-1 : S.textPos]
	}
}
