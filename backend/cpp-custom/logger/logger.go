package logger

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Info writes logs in the color blue with "INFO: " as prefix
var Info = log.New(os.Stdout, "\u001b[34mINFO: \u001B[0m", log.LstdFlags|log.Lshortfile)

// Warning writes logs in the color yellow with "WARNING: " as prefix
var Warning = log.New(os.Stdout, "\u001b[33mWARNING: \u001B[0m", log.LstdFlags|log.Lshortfile)

// Error writes logs in the color red with "ERROR: " as prefix
var Error = log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags|log.Lshortfile)

// Debug writes logs in the color cyan with "DEBUG: " as prefix
var Debug = log.New(os.Stdout, "\u001b[36mDEBUG: \u001B[0m", log.LstdFlags|log.Lshortfile)

// Debug json body write
var DebugJSONBody = struct {
	logger         *log.Logger
	PrintJSONBytes func(r *http.Request)
}{Debug, func(r *http.Request) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error.Println(err)
		return
	}
	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
	var rBodyBuf bytes.Buffer
	_, err = rBodyBuf.ReadFrom(rdr1)
	if err != nil {
		Error.Println(err)
		return
	}
	Debug.Println(rBodyBuf.String())
	r.Body = rdr2
}}
