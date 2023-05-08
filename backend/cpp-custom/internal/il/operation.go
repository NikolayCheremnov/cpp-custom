package il

import "cpp-custom/internal/datatype"

// operation intermediate instructions
const (
	ASSIGN       = iota
	DECLARE_PROC = "DECLARE_PROCEDURE" // declare procedure
	CALL_PROC    = "CALL_PROCEDURE"    // call procedure
	POP          = "POP"               // pop from stack
	PUSH         = "PUSH"              // push to stack
	CAST         = "CAST"              // type converting
)

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
