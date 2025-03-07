package cmd

import (
	"loan-management/infrastructure"
	"log"
)

func Migrate() {
	_, err := infrastructure.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer infrastructure.CloseDB()

	if err := infrastructure.Migrate(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully")
}
