package delivery

import (
	"loan-management/internal/entity"
	"loan-management/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	transactionUsecase *usecase.TransactionUsecase
}

func NewTransactionHandler(transactionUsecase *usecase.TransactionUsecase) *TransactionHandler {
	return &TransactionHandler{transactionUsecase: transactionUsecase}
}

func (h *TransactionHandler) InquiryTransaction(ctx *fiber.Ctx) error {
	loanID, err := strconv.ParseInt(ctx.Query("loan_id"), 10, 64)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	inquiryResult, err := h.transactionUsecase.InquiryTransaction(ctx.Context(), loanID)

	if inquiryResult == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Tagihan tidak ditemukan"})
	}

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": inquiryResult})
}

func (h *TransactionHandler) CreateTransaction(ctx *fiber.Ctx) error {
	var payload entity.CreateTransactionPayload
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	createTransactionPayload := &entity.CreateTransactionPayload{
		LoanID: payload.LoanID,
		Amount: payload.Amount,
	}

	trx, err := h.transactionUsecase.CreateTransaction(ctx.Context(), createTransactionPayload)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if trx == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Tagihan tidak ditemukan"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": trx})

}
