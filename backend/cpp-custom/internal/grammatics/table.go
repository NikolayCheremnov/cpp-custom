// package grammaics
package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"strconv"
	"strings"
)

type LlTable struct {
	Table        map[string]map[string]string
	NonTerminals []string
	Terminals    []string
}

func (table *LlTable) GetTerminalsListByNonTerminal(nonTerminal string) []string {
	row, ok := table.Table[nonTerminal]
	if ok {
		var terminals []string
		for terminal, _ := range row {
			terminals = append(terminals, terminal)
		}
		return terminals
	}
	return nil
}

func (table *LlTable) IsEpsilonRuleExists(nonTerminal string) bool {
	row, ok := table.Table[nonTerminal]
	if ok {
		for _, val := range row {
			if val == "EPSILON" {
				return true
			}
		}
	}
	return false
}

func (table *LlTable) GetFirstNonTerminal() string {
	return table.NonTerminals[0]
}

func (table *LlTable) IsTerminal(sample string) bool {
	for _, terminal := range table.Terminals {
		if sample == terminal {
			return true
		}
	}
	return false
}

func (table *LlTable) IsNonTerminal(sample string) bool {
	for _, nonTerminal := range table.NonTerminals {
		if sample == nonTerminal {
			return true
		}
	}
	return false
}

func (table *LlTable) ConsolePrint() {
	// 1. print terminals
	fmt.Println("Terminals list: " + strings.Join(table.Terminals, " "))
	// 2. print non terminals
	fmt.Println("Non terminals list: " + strings.Join(table.NonTerminals, " "))
	// 3. print values
	for nonTerminal, terminalMap := range table.Table {
		fmt.Println(nonTerminal + ": " + fmt.Sprint(terminalMap))
	}
}

func ReadLLTableFromExcel(excelPath string, tableSheet string, nonTerminalsCellCount int, terminalsCellCount int) (*LlTable, error) {
	// this function can read LL table from excel file to map
	// 1-level keys - non-terminals in rows
	// 2-level-keys - terminals in columns

	llExcelFile, err := excelize.OpenFile(excelPath)
	if err != nil {
		return nil, err
	}

	// 1. read non-terminals list
	var nonTerminalsList []string
	for i := 0; i < nonTerminalsCellCount; i++ {
		val, err := llExcelFile.GetCellValue(tableSheet, "A"+strconv.Itoa(i+2))
		if err != nil {
			return nil, err
		}
		nonTerminalsList = append(nonTerminalsList, val)
	}
	// logger.Info.Print("Read nonTerminalList: " + strings.Join(nonTerminalsList, ", "))

	// 2. read terminals list
	var terminalsKeys []string
	for _, first := range []string{"", "A"} {
		for second := 'A'; second < 'Z' && terminalsCellCount > 0; second++ {
			if first == "" && second == 'A' {
				continue
			}
			terminalsKeys = append(terminalsKeys, first+string(second))
			terminalsCellCount--
		}
	}
	var terminalsList []string
	for _, terminalKey := range terminalsKeys {
		val, err := llExcelFile.GetCellValue(tableSheet, terminalKey+"1")
		if err != nil {
			return nil, err
		}
		terminalsList = append(terminalsList, val)
	}
	// logger.Info.Print("Read terminalList: " + strings.Join(terminalsList, ", "))

	// 3. read table
	llTable := make(map[string]map[string]string)
	for nonTerminalIndex, nonTerminal := range nonTerminalsList {
		llTable[nonTerminal] = make(map[string]string)
		for terminalKeyIndex, terminalKey := range terminalsKeys {
			val, err := llExcelFile.GetCellValue(tableSheet, terminalKey+strconv.Itoa(nonTerminalIndex+2))
			if err != nil {
				return nil, err
			}
			writeValue := func(nonTerminal string, terminal string, value string) {
				if terminal == "a-z" {
					for symb := 'a'; symb <= 'z'; symb++ {
						llTable[nonTerminal][string(symb)] = value
					}
				} else if terminal == "0-9" {
					for symb := '0'; symb <= '9'; symb++ {
						llTable[nonTerminal][string(symb)] = value
					}
				} else {
					llTable[nonTerminal][terminal] = value
				}
			}
			if val != "" {
				if val == "a-z" {
					for symb := 'a'; symb <= 'z'; symb++ {
						writeValue(nonTerminal, string(symb), string(symb))
					}
				} else if val == "0-9" {
					for symb := '0'; symb <= '9'; symb++ {
						writeValue(nonTerminal, string(symb), string(symb))
					}
				} else {
					writeValue(nonTerminal, terminalsList[terminalKeyIndex], val)
				}
			}
		}
	}

	// 4. collect struct
	llTableStruct := new(LlTable)
	// 4.1. add non terminals
	llTableStruct.NonTerminals = nonTerminalsList
	// 4.2 add terminals
	for _, terminal := range terminalsList {
		if terminal == "a-z" {
			for symb := 'a'; symb <= 'z'; symb++ {
				llTableStruct.Terminals = append(llTableStruct.Terminals, string(symb))
			}
		} else if terminal == "0-9" {
			for symb := '0'; symb <= '9'; symb++ {
				llTableStruct.Terminals = append(llTableStruct.Terminals, string(symb))
			}
		} else {
			llTableStruct.Terminals = append(llTableStruct.Terminals, terminal)
		}
	}
	// 4.3 add table
	llTableStruct.Table = llTable

	return llTableStruct, nil
}
