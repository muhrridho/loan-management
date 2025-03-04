package main

import (
	"log"

	"loan-management/infrastructure"
	"loan-management/internal/delivery"
	"loan-management/internal/repository"
	"loan-management/internal/usecase"
	"loan-management/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// if err := infrastructure.Destroy(); err != nil {
	// 	log.Fatalf("Failed to destroy DB: %v", err)
	// }

	db, err := infrastructure.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer infrastructure.CloseDB()

	if err := infrastructure.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// if err := infrastructure.Seed(); err != nil {
	// 	log.Fatalf("Failed to run seeding: %v", err)
	// }

	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := delivery.NewUserHandler(userUsecase)

	loanRepo := repository.NewLoanRepository(db)
	loanUsecase := usecase.NewLoanUsecase(loanRepo)
	loanHandler := delivery.NewLoanHandler(loanUsecase)

	app := fiber.New()

	routes := routes.NewRoutes(app, userHandler, loanHandler)
	routes.SetupRoutes()
	log.Println(app.GetRoutes())

	log.Fatal(app.Listen(":3100"))
}
