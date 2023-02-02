package logger

import (
	"cpp-custom/internal/filesystem"
	"log"
)

var enableLogging bool = false
var fLoggers = make(map[string]*log.Logger)

func Init(loggers map[string]string) error {
	enableLogging = true
	for loggerType, loggerFile := range loggers {
		_, w, err := filesystem.Create("../cpp-custom/data/logs/" + loggerFile + ".log")
		if err != nil {
			return err
		}
		fLoggers[loggerType] = log.New(w, "LOG:", log.LstdFlags|log.Lshortfile)
	}
	return nil
}

func InitWithCustomLogDir(loggers map[string]string, logDir string) error {
	enableLogging = true
	for loggerType, loggerFile := range loggers {
		_, w, err := filesystem.Create(logDir + "/" + loggerFile + ".log")
		if err != nil {
			return err
		}
		fLoggers[loggerType] = log.New(w, "LOG:", log.LstdFlags|log.Lshortfile)
	}
	return nil
}

func Log(logger string, msg string) {
	if enableLogging {
		if fLogger, ok := fLoggers[logger]; ok {
			fLogger.Println(msg)
		} else {
			Error.Println("logger does not exist")
		}
	}
}
