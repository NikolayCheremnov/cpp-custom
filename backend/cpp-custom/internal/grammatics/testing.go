package main

import (
	"cpp-custom/internal/filesystem"
	"io"
)

func main() {
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
}
