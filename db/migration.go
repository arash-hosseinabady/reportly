package db

import (
	"log"
	"reportly/model"
)

// RegisterModels returns all models that should be migrated
func RegisterModels() []interface{} {
	return []interface{}{
		&model.ReportRequest{},
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

		log.Printf("Migrated table: %s", tn)
	}

	log.Println("Migration process completed.")
}
