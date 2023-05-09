package il

import "cpp-custom/internal/datatype"

// operand types
const (
	VARIABLE  = "variable"
	CONSTANT  = "constant"
	OPERATION = "operation"
	PROCEDURE = "procedure"
	TYPE      = "type"
)

type Operand struct {
	Type         string // type of operand
	Identity     string // operand identity
	OperandValue *datatype.Value
}

func (o *Operand) OperandAsString() string {
	out := o.Type + " " + o.Identity
	if o.OperandValue != nil {
		out += " " + o.OperandValue.ValueAsString()
	}
	return out
}
