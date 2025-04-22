package db

import (
	"gorm.io/gorm"
	"log"
)

func GetDb() *gorm.DB {
	return DB
}

// GetTableName returns the table name for a given model instance using GORM's schema parsing.
func GetTableName(model interface{}) string {
	db := GetDb()
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		log.Panic(err)
	}
	return stmt.Schema.Table
}
