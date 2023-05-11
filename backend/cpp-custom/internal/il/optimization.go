package il

import "strconv"

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
		// thereAreDistribution = false
		operations, thereAreDistribution = distributeConstants(operations)
	}
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
	// passCount := 0
	for i := 0; i < len(operations); i++ {
		if operations[i].IlInstruction == STOR && operations[i].RightOperand.Type == CONSTANT {
			// need distribution for this variable below to the nearest stor
			for j := i + 1; j < len(operations); j++ {
				if operations[j].IlInstruction == STOR &&
					operations[j].LeftOperand.Identity == operations[i].LeftOperand.Identity {
					break // stop if meet next stor to this variable
				}

				//if operations[j].IlInstruction == PASS && passCount%2 == 0 {
				//	passCount++
				//	// ahead of the cycle
				//	isCanDistribute := true
				//
				//}

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
	return operations, thereAreDistribution
}
