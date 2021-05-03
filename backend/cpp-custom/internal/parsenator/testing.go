package parsenator

import (
	"io"
	"log"
	"os"
)

func Testing(srsFileName string) (err error) {
	// writers preparing
	var sw, aw io.Writer
	sw, err = os.Create("lexinatorErrors.err")
	if err != nil {
		return err
	}
	aw, err = os.Create("parsenatorErrors.err")
	if err != nil {
		return err
	}

	// analyzer preparation
	A, err := Preparing(srsFileName, sw, aw)
	if err != nil {
		return nil
	}

	// deferred call for panic interception
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred in "+srsFileName+":", err)
		} else {
			log.Println("there are no errors in " + srsFileName)
		}
	}()

	// testing
	err = A.GlobalDescriptions()
	if err != nil {
		return err
	}

	return nil
}
