package logger

import (
	"log"
	"os"
)

var CustomLogger *Logger

// Logger represents a custom logger with separate info and error loggers.
type Logger struct {
	Info  *log.Logger
	Error *log.Logger
}

func init() {
	CustomLogger = &Logger{
		Info:  log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile),
		Error: log.New(os.Stdout, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
