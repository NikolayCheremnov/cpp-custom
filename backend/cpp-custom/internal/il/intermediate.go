package il

import (
	"errors"
	"strconv"
	"strings"
)

const DEFAULT_LIMIT = 10_000

type Intermediate struct {
	Operations      []Operation // array of triads that translate to intermediate code
	OperationsLimit int         // limit of operations (maybe can help in recursion)
}

func NewIntermediateDomain(programSize int) *Intermediate {
	if programSize <= 0 {
		programSize = DEFAULT_LIMIT
	}
	return &Intermediate{[]Operation{}, programSize}
}

func (o *Intermediate) IntermediateAsString() string {
	builder := strings.Builder{}
	for index, value := range o.Operations {
		builder.WriteString(strconv.Itoa(index) + ") " + value.OperationAsString() + "\n")
	}
	return builder.String()
}

func (o *Intermediate) DeclareProcedure(procedureIdentity string) error {
	if len(o.Operations) >= o.OperationsLimit {
		return errors.New("too much instructions")
	}
	operation := Operation{
		IlInstruction: DECLARE_PROC,
		LeftOperand: &Operand{
			Type:         PROCEDURE,
			Identity:     procedureIdentity,
			OperandValue: nil,
		},
		RightOperand: nil,
		Result:       nil,
	}
	o.Operations = append(o.Operations, operation)
	return nil
}

func (o *Intermediate) ExtractArgumentFromStack(argumentIdentity string) error {
	if len(o.Operations) >= o.OperationsLimit {
		return errors.New("too much instructions")
	}
	operation := Operation{
		IlInstruction: POP,
		LeftOperand: &Operand{
			Type:         VARIABLE,
			Identity:     argumentIdentity,
			OperandValue: nil,
		},
		RightOperand: nil,
		Result:       nil,
	}
	o.Operations = append(o.Operations, operation)
	return nil
}

func (o *Intermediate) CallProcedure(procedureIdentity string) error {
	if len(o.Operations) >= o.OperationsLimit {
		return errors.New("too much instructions")
	}
	operation := Operation{
		IlInstruction: CALL_PROC,
		LeftOperand: &Operand{
			Type:         PROCEDURE,
			Identity:     procedureIdentity,
			OperandValue: nil,
		},
		RightOperand: nil,
		Result:       nil,
	}
	o.Operations = append(o.Operations, operation)
	return nil
}

// help methods

func (o *Intermediate) FindLastProcCallOperationIndex() int {
	// find operations for fill left operands
	for i := len(o.Operations) - 1; i >= 0; i-- {
		if o.Operations[i].IlInstruction == CALL_PROC {
			return i
		}
	}
	return -1
}

func (o *Intermediate) MoveOperationToEnd(srcIndex int) {
	item := o.Operations[srcIndex]
	o.Operations = append(o.Operations[:srcIndex], o.Operations[srcIndex+1:]...)
	o.Operations = append(o.Operations, item)
}
