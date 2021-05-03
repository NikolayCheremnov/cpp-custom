package lexinator

import (
	"fmt"
	"io"
	"os"
)

func ScannerTesting(srsFileName string, errFileName string) (err error) {
	// writer preparing
	var w io.Writer
	if errFileName == "" {
		w = os.Stdout
	} else {
		w, err = os.Create(errFileName)
		if err != nil {
			return err
		}
	}

	// scanner preparing
	S := Scanner{sourceModule: "undefined", textPos: 0, line: 0, linePos: 0, writer: w}
	err = S.GetData(srsFileName)
	if err != nil {
		return err
	}
	lexType := -2
	var lexImage string
	for lexType != End {
		lexType, lexImage = S.Scan()
		fmt.Println(lexImage, " type ", lexType)
	}

	fmt.Println("last:")
	lexType, lexImage = S.Scan()
	fmt.Println(lexImage, " type ", lexType)
	return nil
}
