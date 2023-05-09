package il

import "cpp-custom/internal/datatype"

// operation intermediate instructions
const (
	STOR         = "STOR"
	DECLARE_PROC = "DEFP" // declare procedure
	CALL_PROC    = "CALL" // call procedure
	POP          = "POP"  // pop from stack
	PUSH         = "PUSH" // push to stack
	CAST         = "CAST" // type converting
	// operators
	ADD    = "ADD"
	SUB    = "SUB"
	MUL    = "MUL"
	DIV    = "DIV"
	MOD    = "MOD"
	EQ     = "EQ"
	NEQ    = "NEQ"
	MORE   = "MORE"
	LESS   = "LESS"
	MOREEQ = "MOREEQ"
	LESSEQ = "LESSEQ"
	//
	PASS = "PASS" // empty operation
	//
	JMP  = "JMP"  // no conditional jump to operation
	JMPF = "JMPF" // conditional jump if right operand is false
)

func GetOperatorByLexeme(operatorLex string) string {
	switch operatorLex {
	case "+":
		return ADD
	case "-":
		return SUB
	case "*":
		return MUL
	case "/":
		return DIV
	case "%":
		return MOD
	case "==":
		return EQ
	case "!=":
		return NEQ
	case ">":
		return MORE
	case "<":
		return LESS
	case ">=":
		return MOREEQ
	case "<=":
		return LESSEQ
	default:
		panic("Invalid operator '" + operatorLex + "'")
	}
}

// Operation - implementation of triad values
type Operation struct {
	IlInstruction string
	LeftOperand   *Operand
	RightOperand  *Operand
	Result        *datatype.Value
}

func (o *Operation) OperationAsString() string {
	out := o.IlInstruction
	if o.LeftOperand != nil {
		out += " (" + o.LeftOperand.OperandAsString() + ")"
	} else {
		out += " (NONE)"
	}
	if o.RightOperand != nil {
		out += " (" + o.RightOperand.OperandAsString() + ")"
	} else {
		out += " (NONE)"
	}
	if o.Result != nil {
		out += " -> " + o.Result.ValueAsString()
	}
	return out
}
