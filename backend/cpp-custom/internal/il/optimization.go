package il

import (
	"fmt"
	"strconv"
)

func MakeIntermediateOptimization(operations []Operation) ([]Operation, error) {
	thereAreDistribution := true
	for thereAreDistribution {
		thereAreCalculations := true
		for thereAreCalculations {
			var err error
			operations, thereAreCalculations, err = constantsConvolution(operations)
			if err != nil {
				return operations, err
			}
			// 2. apply operation-operands with known result
			operations = applyOperandsWithResult(operations)
			// 3. remove not used calculation instructions
			operations = removeNotUsedCalculationOperations(operations)
		}
		operations = removeUnnecessaryTypeCasts(operations)
		operations = removeUnnecessaryStors(operations)
		// thereAreDistribution = false
		operations, thereAreDistribution = distributeConstants(operations)
	}
	operations = removeUnreachableCode(operations)
	return operations, nil
}

func constantsConvolution(operations []Operation) ([]Operation, bool, error) {
	thereAreCalculations := false
	for i := 0; i < len(operations); i++ {
		isCalculated, err := operations[i].calculateResultIfBinaryWithConstants()
		if isCalculated {
			thereAreCalculations = true
		}
		if err != nil {
			return operations, false, err
		}
	}

	return operations, thereAreCalculations, nil
}

func applyOperandsWithResult(operations []Operation) []Operation {
	for i := 0; i < len(operations); i++ {
		if operations[i].LeftOperand != nil && operations[i].LeftOperand.Type == OPERATION {
			operationRef, _ := strconv.Atoi(operations[i].LeftOperand.Identity)
			if operations[operationRef].Result != nil {
				operations[i].LeftOperand = &Operand{
					Type:         CONSTANT,
					Identity:     strconv.FormatInt(operations[operationRef].Result.DataAsInt, 10),
					OperandValue: operations[operationRef].Result,
				}
			}
		}
		if operations[i].RightOperand != nil && operations[i].RightOperand.Type == OPERATION {
			operationRef, _ := strconv.Atoi(operations[i].RightOperand.Identity)
			if operations[operationRef].Result != nil {
				operations[i].RightOperand = &Operand{
					Type:         CONSTANT,
					Identity:     strconv.FormatInt(operations[operationRef].Result.DataAsInt, 10),
					OperandValue: operations[operationRef].Result,
				}
			}
		}
	}
	return operations
}

func removeNotUsedCalculationOperations(operations []Operation) []Operation {
	for i := 0; i < len(operations); {
		if isCalculationOperation(operations[i].IlInstruction) {
			isNotUsed := true
			for j := i + 1; j < len(operations); j++ {
				if operations[j].LeftOperand != nil && operations[j].LeftOperand.Type == OPERATION &&
					operations[j].LeftOperand.Identity == strconv.Itoa(i) ||
					operations[j].RightOperand != nil && operations[j].RightOperand.Type == OPERATION &&
						operations[j].RightOperand.Identity == strconv.Itoa(i) {
					isNotUsed = false
					break
				}
			}
			if isNotUsed {
				operations = removeOperation(operations, i)
			} else {
				i++
			}
		} else {
			i++
		}
	}
	return operations
}

func removeOperation(operations []Operation, removed int) []Operation {
	operations = append(operations[:removed], operations[removed+1:]...)
	// need to decrease refs
	for i := removed; i < len(operations); i++ {
		if operations[i].IlInstruction != JMP && operations[i].IlInstruction != JMPF {
			if operations[i].LeftOperand != nil && operations[i].LeftOperand.Type == OPERATION {
				ref, _ := strconv.Atoi(operations[i].LeftOperand.Identity)
				if ref > removed {
					operations[i].LeftOperand.Identity = strconv.Itoa(ref - 1)
				}
			}
			if operations[i].RightOperand != nil && operations[i].RightOperand.Type == OPERATION {
				ref, _ := strconv.Atoi(operations[i].RightOperand.Identity)
				if ref > removed {
					operations[i].RightOperand.Identity = strconv.Itoa(ref - 1)
				}
			}
		}
	}
	// need to decrease JMP and JMPF operations
	for i := 0; i < len(operations); i++ {
		if operations[i].IlInstruction == JMP || operations[i].IlInstruction == JMPF {
			ref, _ := strconv.Atoi(operations[i].LeftOperand.Identity)
			if ref >= removed {
				operations[i].LeftOperand.Identity = strconv.Itoa(ref - 1)
			}
		}
		if operations[i].IlInstruction == JMPF {
			ref, _ := strconv.Atoi(operations[i].RightOperand.Identity)
			if ref >= removed {
				operations[i].RightOperand.Identity = strconv.Itoa(ref - 1)
			}
		}
	}
	return operations
}

func removeUnnecessaryTypeCasts(operations []Operation) []Operation {
	for i := 0; i < len(operations); {
		if operations[i].IlInstruction == CAST && operations[i].LeftOperand.Type == CONSTANT &&
			operations[i].RightOperand.Identity == "int" {
			// replace left operand if next is push, replace right operand if next stor for next operation
			if operations[i+1].IlInstruction == PUSH {
				operations[i+1].LeftOperand = operations[i].LeftOperand
			} else {
				// else stor
				operations[i+1].RightOperand = operations[i].LeftOperand
			}
			operations = removeOperation(operations, i)
		} else {
			i++
		}
	}
	return operations
}

func distributeConstants(operations []Operation) ([]Operation, bool) {
	thereAreDistribution := false

	for i := 0; i < len(operations); i++ {
		if operations[i].IlInstruction == STOR && operations[i].RightOperand.Type == CONSTANT {
			// need distribution for this variable below to the nearest stor
			for j := i + 1; j < len(operations); j++ {
				fmt.Printf("i = %d, j = %d\n", i, j)
				if i == 78 && j == 79 {
					a := 1
					a++
				}
				if operations[j].IlInstruction == STOR &&
					operations[j].LeftOperand.Identity == operations[i].LeftOperand.Identity {
					break // stop if meet next stor to this variable
				} else if operations[j].IlInstruction == PASS && operations[j-1].IlInstruction != JMP {
					// start cycle from j - need check that there are no non-constant stor to distribution variable
					isCanDistribute := true
					for k := j + 1; operations[k].IlInstruction != JMP &&
						(operations[k].LeftOperand == nil || operations[k].LeftOperand.Identity != strconv.Itoa(j)); k++ {
						if operations[k].IlInstruction == STOR &&
							operations[k].LeftOperand.Identity == operations[i].LeftOperand.Identity &&
							operations[k].RightOperand.Type != CONSTANT {
							isCanDistribute = false
							break
						}
					}
					if !isCanDistribute {
						break // can`t distribute then changed in cycle without constant
					}
				} else {
					if operations[j].LeftOperand != nil && operations[j].LeftOperand.Identity == operations[i].LeftOperand.Identity {
						operations[j].LeftOperand = operations[i].RightOperand
						thereAreDistribution = true
					}
					if operations[j].RightOperand != nil && operations[j].RightOperand.Identity == operations[i].LeftOperand.Identity {
						operations[j].RightOperand = operations[i].RightOperand
						thereAreDistribution = true
					}
				}
			}
		}
	}
	return operations, thereAreDistribution
}

func removeUnnecessaryStors(operations []Operation) []Operation {
	for i := 0; i < len(operations); {
		if operations[i].IlInstruction == STOR {
			// find nearest stor and shoul only calculations operations between (expression)
			isUnnecessaryStor := false
			for j := i + 1; j < len(operations); j++ {
				if operations[j].IlInstruction == STOR && operations[j].LeftOperand.Identity ==
					operations[i].LeftOperand.Identity {
					isUnnecessaryStor = true
				} else if !isCalculationOperation(operations[j].IlInstruction) {
					break
				}
			}
			if isUnnecessaryStor {
				operations = removeOperation(operations, i)
			} else {
				i++
			}
		} else {
			i++
		}
	}
	return operations
}

func removeUnreachableCode(operations []Operation) []Operation {
	for i := 0; i < len(operations); {
		if operations[i].IlInstruction == JMPF &&
			operations[i-1].LeftOperand.Type == CONSTANT && !operations[i-1].LeftOperand.OperandValue.DataAsBool {
			// always false -> need remove cycle
			cycleStart := i - 2
			var j int
			for j = i - 1; operations[j].IlInstruction != JMP && operations[j].LeftOperand.Identity != strconv.Itoa(cycleStart); j++ {
			}
			cycleEnd := j + 1
			for k := 0; k < cycleEnd-cycleStart+1; k++ {
				operations = removeOperation(operations, cycleStart)
			}
		} else {
			i++
		}
	}
	return operations
}
