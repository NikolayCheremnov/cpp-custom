package polis

func precedence(operator string) int {
	switch operator {
	case "*", "/", "%":
		return 2
	case "+", "-":
		return 1
	case "==", "!=", ">", "<", ">=", "<=":
		return 0
	default:
		return -1
	}
}

func IsOperator(token string) bool {
	operators := []string{"+", "-", "*", "/", "%", "==", "!=", ">", "<", ">=", "<="}
	for _, op := range operators {
		if op == token {
			return true
		}
	}
	return false
}

func ConvertToRPN(tokens []string) []string {
	stack := []string{}
	output := []string{}

	for _, token := range tokens {
		if IsOperator(token) {
			for len(stack) > 0 && IsOperator(stack[len(stack)-1]) && precedence(token) <= precedence(stack[len(stack)-1]) {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		} else if token == "(" {
			stack = append(stack, token)
		} else if token == ")" {
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) > 0 && stack[len(stack)-1] == "(" {
				stack = stack[:len(stack)-1] // Удаление "(" из стека
			}
		} else {
			output = append(output, token)
		}
	}

	for len(stack) > 0 {
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return output
}

func AddUnaryMinus(tokens []string) []string {
	result := []string{}

	isPass := false
	for i, token := range tokens {
		if isPass {
			isPass = false
			continue
		}
		if token == "-" && (i == 0 || IsOperator(tokens[i-1]) || tokens[i-1] == "(") {
			result = append(result, "(")
			result = append(result, "0")
			result = append(result, token)
			result = append(result, tokens[i+1])
			result = append(result, ")")
			isPass = true
		} else {
			result = append(result, token)
		}

	}

	return result
}
