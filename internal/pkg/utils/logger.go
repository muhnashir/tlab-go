package utils

import (
	"log"
	"os"
)

var (
	InfoLogger  = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)

// LogInfo prints info messages
func LogInfo(v ...interface{}) {
	InfoLogger.Println(v...)
}

// LogInfof prints formatted info messages
func LogInfof(format string, v ...interface{}) {
	InfoLogger.Printf(format, v...)
}

// LogError prints error messages
func LogError(v ...interface{}) {
	ErrorLogger.Println(v...)
}

// LogErrorf prints formatted error messages
func LogErrorf(format string, v ...interface{}) {
	ErrorLogger.Printf(format, v...)
}

// LogCheckError checks if err is not nil, then logs it as error
func LogCheckError(err error, message string) {
	if err != nil {
		ErrorLogger.Printf("%s: %v", message, err)
	}
}
