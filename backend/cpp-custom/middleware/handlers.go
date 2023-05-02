package middleware

import (
	"cpp-custom/internal/filesystem"
	"cpp-custom/internal/ll1"
	"cpp-custom/internal/parsenator"
	"cpp-custom/internal/semanthoid"
	"cpp-custom/logger"
	"cpp-custom/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var Port string

func Ping(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(basic_response{Message: "ping OK, port " + Port}) // TODO: added runserver configurations
	logger.Info.Println("ping OK, port "+Port, r.RemoteAddr)
}

func ProcessCodeByLl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// read srs code
	var srsText models.SrsText
	err := json.NewDecoder(r.Body).Decode(&srsText)
	if err_handling(err, "can`t decode source text", w) != nil {
		return
	}

	// write srs code to tmp file
	srsFilePath, srsFile, err := filesystem.Create("../cpp-custom/data/tmp/srs.cpp")
	if err_handling(err, "can`t create tmp srs file", w) != nil {
		return
	}
	_, err = srsFile.WriteString(srsText.Text)
	if err_handling(err, "problem with tmp files", w) != nil {
		return
	}
	if er := srsFile.Close(); err_handling(er, "problem with tmp files during Close()", w) != nil {
		return
	}

	// error writers preparation
	var sw, cw io.Writer
	swPath, sw, err := filesystem.Create("../cpp-custom/data/tmp/lexinatorErrors.err")
	if err_handling(err, "problem with tmp files", w) != nil {
		return
	}
	cwPath, cw, err := filesystem.Create("../cpp-custom/data/tmp/checkerErrors.err")
	if err_handling(err, "problem with tmp files", w) != nil {
		return
	}

	// checker preparation
	checker, err := ll1.CreateLlChecker(srsFilePath, sw, cw)
	if err_handling(err, "can`t prepare checker", w) != nil {
		return
	}

	// deferred call for panic interception
	defer func() {
		lexinatorErrors, er := filesystem.ReadFileToString(swPath)
		if err_handling(er, "can`t read lexinator errors", w) != nil {
			return
		}
		checkerErrors, er := filesystem.ReadFileToString(cwPath)
		if err_handling(err, "can`t read checker errors", w) != nil {
			return
		}
		var Message string
		if panicMsg := recover(); panicMsg != nil {
			Message = "panic occurred: " + fmt.Sprint(panicMsg) + " check errors notes"
		} else {
			Message = checker.TreeToString()
			lexinatorErrors = "there are no errors"
			checkerErrors = "there are no errors"
		}
		logger.Info.Println(Message)
		json.NewEncoder(w).Encode(CheckerResponse{
			Message:         Message,
			LexinatorErrors: lexinatorErrors,
			CheckerErrors:   checkerErrors,
		})
	}()

	// run processing
	checker.MakeLkAnalyze()
}

func CheckForErrors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// read srs code
	var srsText models.SrsText
	err := json.NewDecoder(r.Body).Decode(&srsText)
	if err_handling(err, "can`t decode source text", w) != nil {
		return
	}

	// write srs code to tmp file
	srsFilePath, srsFile, err := filesystem.Create("../cpp-custom/data/tmp/srs.cpp")
	if err_handling(err, "can`t create tmp srs file", w) != nil {
		return
	}
	_, err = srsFile.WriteString(srsText.Text)
	if err_handling(err, "problem with tmp files", w) != nil {
		return
	}
	srsFile.Close()

	// error writers preparation
	var sw, aw io.Writer
	swPath, sw, err := filesystem.Create("../cpp-custom/data/tmp/lexinatorErrors.err")
	if err_handling(err, "problem with tmp files", w) != nil {
		return
	}
	awPath, aw, err := filesystem.Create("../cpp-custom/data/tmp/parsenatorErrors.err")
	if err_handling(err, "problem with tmp files", w) != nil {
		return
	}

	// analyzer preparation
	A, err := parsenator.Preparing(srsFilePath, sw, aw)
	if err_handling(err, "can`t prepare analyzer", w) != nil {
		return
	}

	// deferred call for panic interception
	defer func() {
		LexinatorErrors, err := filesystem.ReadFileToString(swPath)
		if err_handling(err, "can`t read lexinator errors", w) != nil {
			return
		}
		ParsenatorErrors, err := filesystem.ReadFileToString(awPath)
		if err_handling(err, "can`t read parsenator errors", w) != nil {
			return
		}
		var Message string
		if err := recover(); err != nil {
			Message = "panic occurred: " + fmt.Sprint(err) + " check errors notes"
		} else {
			Message = semanthoid.TreeToString()
			ParsenatorErrors = "there are no errors"
			LexinatorErrors = "there are no errors"
		}
		logger.Info.Println(Message)
		json.NewEncoder(w).Encode(analyze_response{
			Message:          Message,
			LexinatorErrors:  LexinatorErrors,
			ParsenatorErrors: ParsenatorErrors,
		})
	}()

	// testing
	err = A.GlobalDescriptions()
	if err_handling(err, "bad analyze running", w) != nil {
		return
	}
}
