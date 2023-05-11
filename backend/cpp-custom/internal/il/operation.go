package il

import (
	"cpp-custom/internal/datatype"
	"errors"
)

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

func isCalculationOperation(operation string) bool {
	return operation == ADD || operation == SUB || operation == MUL || operation == DIV || operation == MOD ||
		operation == EQ || operation == NEQ || operation == MORE || operation == LESS ||
		operation == MOREEQ || operation == LESSEQ
}

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

// calculateResultIfBinaryWithConstants calculates result and store to Result
// only for calculations operators and constant operands
// never do if not calculation operator or have not constant operands
func (o *Operation) calculateResultIfBinaryWithConstants() (bool, error) {
	isCalculated := false
	if isCalculationOperation(o.IlInstruction) {
		if o.LeftOperand == nil || o.RightOperand == nil {
			return false, errors.New("calculation operator must have two operands")
		}
		if o.LeftOperand.Type == CONSTANT && o.RightOperand.Type == CONSTANT {
			var err error
			o.Result, err = applyCalculationOperator(o.IlInstruction, o.LeftOperand.OperandValue, o.RightOperand.OperandValue)
			if err != nil {
				return false, err
			}
			isCalculated = true
		}
	}
	return isCalculated, nil
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

func applyCalculationOperator(operator string, o1 *datatype.Value, o2 *datatype.Value) (*datatype.Value, error) {
	var result datatype.Value
	switch operator {
	case ADD:
		result = datatype.NewIntDataValue(o1.DataAsInt + o2.DataAsInt)
		break
	case SUB:
		result = datatype.NewIntDataValue(o1.DataAsInt - o2.DataAsInt)
		break
	case MUL:
		result = datatype.NewIntDataValue(o1.DataAsInt * o2.DataAsInt)
		break
	case DIV:
		result = datatype.NewIntDataValue(o1.DataAsInt / o2.DataAsInt)
		break
	case MOD:
		result = datatype.NewIntDataValue(o1.DataAsInt % o2.DataAsInt)
		break
	case EQ:
		result = datatype.NewBoolDataValue(o1.DataAsInt == o2.DataAsInt)
		break
	case NEQ:
		result = datatype.NewBoolDataValue(o1.DataAsInt != o2.DataAsInt)
		break
	case MORE:
		result = datatype.NewBoolDataValue(o1.DataAsInt > o2.DataAsInt)
		break
	case LESS:
		result = datatype.NewBoolDataValue(o1.DataAsInt < o2.DataAsInt)
		break
	case MOREEQ:
		result = datatype.NewBoolDataValue(o1.DataAsInt >= o2.DataAsInt)
		break
	case LESSEQ:
		result = datatype.NewBoolDataValue(o1.DataAsInt <= o2.DataAsInt)
		break
	default:
		return nil, errors.New("can`t apply not calculation operator")
	}
	return &result, nil
}
