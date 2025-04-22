package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

// InitDB open database connection
func InitDB() {
	var err error
	var dsn string
	var dialector gorm.Dialector

	switch os.Getenv("DB_CONNECTION") {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASS"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)
		dialector = mysql.Open(dsn)
	case "postgres":
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASS"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
		)
		dialector = postgres.Open(dsn)
	default:
		log.Panic("undefined database connection")
	}

	DB, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}

	log.Println("Connected to database")
}

// CloseDB ensures the connection is closed when the application exits
func CloseDB() {
	if DB != nil {
		sqlDb, err := DB.DB()
		if err != nil {
			log.Panic(err)
		}

		err = sqlDb.Close()
		if err != nil {
			log.Panic(err)
		}
	}
}
