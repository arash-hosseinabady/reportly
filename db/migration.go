package db

import (
	"gorm.io/gorm"
	"log"
	"reportly/model"
)

// RegisterModels returns all models that should be migrated
func RegisterModels() []interface{} {
	return []interface{}{
		&model.ReportRequest{},
	}
}

func dropUnusedColumns(db *gorm.DB, model interface{}) {
	stmt := &gorm.Statement{DB: db}
	_ = stmt.Parse(model)

	tableName := stmt.Schema.Table

	// get all current fields
	columns, _ := db.Migrator().ColumnTypes(model)

	for _, col := range columns {
		colName := col.Name()
		// remove field that not exists in a model
		if stmt.Schema.LookUpField(colName) == nil {
			log.Printf("Dropping unused column: %s.%s", tableName, colName)
			_ = db.Migrator().DropColumn(model, colName)
		}
	}
}

func RunMigration() {
	db := GetDb()

	log.Println("Running migrations per model")

	for _, m := range RegisterModels() {
		tn := GetTableName(m)

		log.Printf("Migrating table: %s", tn)

		if err := db.AutoMigrate(m); err != nil {
			log.Printf("Migration failed for table '%s': %v", tn, err)
			continue
		}

		dropUnusedColumns(db, m)

		log.Printf("Migrated table: %s", tn)
	}

	log.Println("Migration process completed.")
}
