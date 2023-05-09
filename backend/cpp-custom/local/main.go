package main

import "cpp-custom/internal/ll1"

func main() {
	// 1. convert long grammar text to short grammar text
	srcLongGrammarTextPath := "D:\\CherepNick\\ASTU\\magistracy\\1_course\\2_semester\\CD\\лр\\cpp-custom\\backend\\cpp-custom\\internal\\ll1\\specifications\\longGrammarWithOperationalSymbols.gr"
	dstShortGrammarTextPath := "D:\\CherepNick\\ASTU\\magistracy\\1_course\\2_semester\\CD\\лр\\cpp-custom\\backend\\cpp-custom\\internal\\ll1\\specifications\\shortGrammarWithOperationalSymbols.gr"
	nonTerminalsJsonPath := "D:\\CherepNick\\ASTU\\magistracy\\1_course\\2_semester\\CD\\лр\\cpp-custom\\backend\\cpp-custom\\internal\\ll1\\specifications\\non-terminals.json"
	ll1.FConvertSourceGrammarText(srcLongGrammarTextPath, dstShortGrammarTextPath, true, nonTerminalsJsonPath)

	// test POLIS
	//tokens := []string{"(", "-", "2", "+", "2", ")", "*", "-2"}
	//tokens := []string{"2", "*", "-", "2"}
	//
	//tokensWithUnaryMinus := polis.AddUnaryMinus(tokens)
	//for _, token := range tokensWithUnaryMinus {
	//	fmt.Print(token + " ")
	//}
	//fmt.Println()
	//rpn := polis.ConvertToRPN(tokensWithUnaryMinus)
	//fmt.Println(strings.Join(rpn, " ")) // Вывод: 2 2 -2 * +
}
