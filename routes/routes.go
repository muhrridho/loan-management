package routes

import (
	"loan-management/internal/delivery"

	"github.com/gofiber/fiber/v2"
)

type Routes struct {
	app                *fiber.App
	userHandler        *delivery.UserHandler
	paymentHandler     *delivery.PaymentHandler
	loanHandler        *delivery.LoanHandler
	transactionHandler *delivery.TransactionHandler
}

func NewRoutes(
	app *fiber.App,
	userHandler *delivery.UserHandler,
	paymentHandler *delivery.PaymentHandler,
	loanHandler *delivery.LoanHandler,
	transactionHandler *delivery.TransactionHandler,
) *Routes {
	return &Routes{
		app:                app,
		userHandler:        userHandler,
		paymentHandler:     paymentHandler,
		loanHandler:        loanHandler,
		transactionHandler: transactionHandler,
	}
}

func (r *Routes) SetupRoutes() {

	api := r.app.Group("/api")

	// Users Group
	users := api.Group("/users")
	users.Get("/", func(ctx *fiber.Ctx) error { return r.userHandler.GetAllUsers(ctx) })
	users.Get("/:id", func(ctx *fiber.Ctx) error { return r.userHandler.GetUserByID(ctx) })
	users.Get("/:id/delinquent-status", func(ctx *fiber.Ctx) error { return r.userHandler.CheckUserDelinquentStatus(ctx) })
	users.Post("/register", func(ctx *fiber.Ctx) error { return r.userHandler.RegisterUser(ctx) })

	// Payment Group
	payments := api.Group("/payments")
	payments.Get("/", func(ctx *fiber.Ctx) error { return r.paymentHandler.GetAllPayments(ctx) })

	// Loans Group
	loans := api.Group("/loans")
	loans.Get("/", func(ctx *fiber.Ctx) error { return r.loanHandler.GetAllLoans(ctx) })
	loans.Get("/:id", func(ctx *fiber.Ctx) error { return r.loanHandler.GetLoanByID(ctx) })
	loans.Post("/create", func(ctx *fiber.Ctx) error { return r.loanHandler.CreateLoan(ctx) })

	// Transaction Group
	trx := api.Group("/transaction")
	trx.Get("/inquiry", func(ctx *fiber.Ctx) error { return r.transactionHandler.InquiryTransaction(ctx) })
	trx.Post("/create", func(ctx *fiber.Ctx) error { return r.transactionHandler.CreateTransaction(ctx) })
}
