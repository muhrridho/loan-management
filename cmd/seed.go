package cmd

import (
	"loan-management/infrastructure"
	"log"
)

func Seed() {
	_, err := infrastructure.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer infrastructure.CloseDB()

	if err := infrastructure.Seed(); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}

	log.Println("Seeding completed successfully")
}
