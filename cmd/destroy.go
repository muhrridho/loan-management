package cmd

import (
	"loan-management/infrastructure"
	"log"
)

func Destroy() {
	if err := infrastructure.Destroy(); err != nil {
		log.Fatalf("Destroy failed: %v", err)
	}

	log.Println("Database destroyed successfully")
}
