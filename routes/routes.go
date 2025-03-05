package routes

import (
	"loan-management/internal/delivery"

	"github.com/gofiber/fiber/v2"
)

type Routes struct {
	app            *fiber.App
	userHandler    *delivery.UserHandler
	paymentHandler *delivery.PaymentHandler
	loanHandler    *delivery.LoanHandler
}

func NewRoutes(app *fiber.App, userHandler *delivery.UserHandler, paymentHandler *delivery.PaymentHandler, loanHandler *delivery.LoanHandler) *Routes {
	return &Routes{
		app:            app,
		userHandler:    userHandler,
		paymentHandler: paymentHandler,
		loanHandler:    loanHandler,
	}
}

func (r *Routes) SetupRoutes() {

	api := r.app.Group("/api")

	// Users Group
	users := api.Group("/users")
	users.Get("/", func(ctx *fiber.Ctx) error { return r.userHandler.GetAll(ctx) })
	users.Get("/:id", func(ctx *fiber.Ctx) error { return r.userHandler.GetByID(ctx) })
	users.Post("/register", func(ctx *fiber.Ctx) error { return r.userHandler.RegisterUser(ctx) })

	// Payment Group
	payments := api.Group("/payments")
	payments.Get("/", func(ctx *fiber.Ctx) error { return r.paymentHandler.GetAllPayments(ctx) })
	// payments.Get("/:loan-id", func(ctx *fiber.Ctx) error { return r.paymentHandler.GetAllPayments(ctx) })

	// Loans Group
	loans := api.Group("/loans")
	loans.Get("/", func(ctx *fiber.Ctx) error { return r.loanHandler.GetAllLoans(ctx) })
	loans.Get("/:id", func(ctx *fiber.Ctx) error { return r.loanHandler.GetLoanByID(ctx) })
	loans.Post("/create", func(ctx *fiber.Ctx) error { return r.loanHandler.CreateLoan(ctx) })
}
