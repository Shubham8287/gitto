package db

import (
	"log"

	f "gitto/features"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const DB_FILE = "assistant.db"

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("assistant.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Migrate the schema (automatically creates the table if it does not exist)
	err = db.AutoMigrate(&f.UserInfo{}, &f.Todo{})
	if err != nil {
		log.Fatalf("Error migrating database schema: %v", err)
	}

	return db
}
