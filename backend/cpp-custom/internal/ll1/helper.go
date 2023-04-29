// package grammaics
package ll1

import (
	"encoding/json"
	"os"
	"strings"
)

func FConvertSourceGrammarText(sourcePath string, compressedPath string, isOnCompression bool) {
	src, dst := "", ""
	if isOnCompression {
		src, dst = sourcePath, compressedPath
	} else {
		src, dst = compressedPath, sourcePath
	}

	file, err := os.ReadFile(src)
	if err != nil {
		panic(err)
	}
	grammarText := string(file)
	compressedText := convertSourceGrammarText(grammarText, isOnCompression)
	err = os.WriteFile(dst, []byte(compressedText), 0644)
	if err != nil {
		panic(err)
	}
}

func convertSourceGrammarText(grammarText string, isOnCompression bool) string {
	for longNt, shortNt := range readGrammarMap() {
		if isOnCompression {
			grammarText = strings.ReplaceAll(grammarText, longNt, shortNt)
		} else {
			grammarText = strings.ReplaceAll(grammarText, shortNt, longNt)
		}
	}
	return grammarText
}

func readGrammarMap() map[string]string {
	file, err := os.ReadFile("./specifications/non-terminals.json")
	if err != nil {
		panic(err)
	}
	ntMap := map[string]string{}
	err = json.Unmarshal([]byte(file), &ntMap)
	if err != nil {
		panic(err)
	}
	return ntMap
}
