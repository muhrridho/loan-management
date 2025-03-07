package main

import (
	"fmt"
	"log"
	"os"

	"loan-management/cmd"
	"loan-management/infrastructure"
	"loan-management/internal/delivery"
	"loan-management/internal/repository"
	"loan-management/internal/usecase"
	"loan-management/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	if len(os.Args) > 1 {
		command := os.Args[1]

		switch command {
		case "migrate":
			cmd.Migrate()
		case "seed":
			cmd.Seed()
		case "destroy":
			cmd.Destroy()
		default:
			fmt.Println("Unknown command:", command)
			fmt.Println("Usage: app [migrate|seed|destroy]")
			os.Exit(1)
		}
		return
	}

	db, err := infrastructure.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer infrastructure.CloseDB()

	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := delivery.NewUserHandler(userUsecase)

	paymentRepo := repository.NewPaymentRepository(db)
	paymentUsecase := usecase.NewPaymentUsecase(paymentRepo)
	paymentHandler := delivery.NewPaymentHandler(paymentUsecase)

	loanRepo := repository.NewLoanRepository(db)
	loanUsecase := usecase.NewLoanUsecase(loanRepo, userUsecase, paymentUsecase)
	loanHandler := delivery.NewLoanHandler(loanUsecase)

	userUsecase.InjectDependencies(loanUsecase)

	transactionRepo := repository.NewTransactionRepository(db)
	transactionUsecase := usecase.NewTransactionUsecase(transactionRepo, loanUsecase, paymentUsecase)
	transactionHandler := delivery.NewTransactionHandler(transactionUsecase)

	app := fiber.New()

	routes := routes.NewRoutes(app, userHandler, paymentHandler, loanHandler, transactionHandler)
	routes.SetupRoutes()
	log.Println(app.GetRoutes())

	log.Fatal(app.Listen(":3100"))
}
