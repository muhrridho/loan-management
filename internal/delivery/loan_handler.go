package delivery

import (
	"fmt"
	"loan-management/internal/entity"
	"loan-management/internal/usecase"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type LoanHandler struct {
	loanUsecase *usecase.LoanUsecase
}

func NewLoanHandler(loanUsecase *usecase.LoanUsecase) *LoanHandler {
	return &LoanHandler{loanUsecase: loanUsecase}
}

func (h *LoanHandler) CreateLoan(ctx *fiber.Ctx) error {
	var payload entity.CreateLoanPayload
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	fmt.Println(ctx.BodyParser(&payload))

	loan := entity.Loan{
		UserID:           payload.UserID,
		Amount:           payload.Amount,
		Interest:         payload.Interest,
		InterestType:     payload.InterestType,
		Tenure:           payload.Tenure,
		TenureType:       payload.TenureType,
		Status:           entity.LoanStatusActive,
		CreatedAt:        time.Now(),
		BillingStartDate: payload.BillingStartDate,
	}

	if err := h.loanUsecase.CreateLoan(ctx.Context(), &loan); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"data": loan})
}

func (h *LoanHandler) GetAllLoans(ctx *fiber.Ctx) error {
	loans, err := h.loanUsecase.GetAllLoans(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if loans == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"data": []entity.Loan{}})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": loans})
}

func (h *LoanHandler) GetLoanByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	loan, err := h.loanUsecase.GetLoanByID(ctx.Context(), id)

	if loan == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Loan not found"})
	}

	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": loan})
}
