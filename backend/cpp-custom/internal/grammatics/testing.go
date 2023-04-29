package main

import (
	"cpp-custom/internal/filesystem"
	"fmt"
	"io"
)

func main() {

	FConvertSourceGrammarText("./specifications/longGrammar.gr", "./specifications/shortGrammar.gr", true)

	// error writers preparation
	var sw, cw io.Writer
	_, sw, err := filesystem.Create("./tdata/lexinatorErrors.err")
	if err != nil {
		panic(err)
	}
	_, cw, err = filesystem.Create("./tdata/llErrors.err")
	if err != nil {
		panic(err)
	}

	//  checker preparation
	checker, err := CreateLlChecker("./tdata/src.cpp", sw, cw)
	if err != nil {
		panic(err)
	}

	// check
	checker.MakeLkAnalyze()
	stackTrace := checker.stackToString()
	fmt.Println(stackTrace)
}
