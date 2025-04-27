package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"reportly/db"
	appLog "reportly/helper"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func initLog() *os.File {
	// Initialize logging
	logFile := appLog.InitLogFile()
	//defer func(logFile *os.File) {
	//	err := logFile.Close()
	//	if err != nil {
	//		log.Panic(err)
	//	}
	//}(logFile) // Close log file when done

	// Redirect panics to the log file
	defer func() {
		if err := recover(); err != nil {
			log.Panic(fmt.Errorf("%v", err))
		}
	}()

	return logFile
}

func main() {
	loadEnv()
	logFile := initLog()

	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			log.Panic(err)
		}
	}(logFile) // Close log file when done

	db.InitDB()
	db.RunMigration()

	Run()

	return
}
