package helper

import (
	"fmt"
	"log"
	"os"
	"time"
)

// InitLogFile Initialize log file based on the current date
func InitLogFile() *os.File {
	// Get current date (YYYY-MM-DD)
	currentDate := time.Now().Format("2006-01-02")

	// Define log file name
	logFileName := fmt.Sprintf("logs/%s.log", currentDate)

	// Create logs directory if not exists
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", 0755)
		if err != nil {
			panic(err)
		}
	}

	// Open log file for appending (create if not exists)
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Set log output to file
	log.SetOutput(logFile)

	// Log panic messages
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	return logFile
}
