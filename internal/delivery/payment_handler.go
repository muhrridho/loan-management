package delivery

import (
	"loan-management/internal/entity"
	"loan-management/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PaymentHandler struct {
	paymentUsecase *usecase.PaymentUsecase
}

func NewPaymentHandler(paymentUsecase *usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{paymentUsecase: paymentUsecase}
}

func (h *PaymentHandler) GetAllPayments(ctx *fiber.Ctx) error {
	var status *entity.PaymentStatus

	if ctx.Params("status") != "" {
		_status, err := strconv.ParseInt(ctx.Params("status"), 10, 8)
		if err != nil {
			return err
		}

		temp := entity.PaymentStatus(_status)
		status = &temp
	}

	payments, err := h.paymentUsecase.GetAllPayments(ctx.Context(), status)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if payments == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"data": []entity.Payment{}})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": payments})
}

// func (h *PaymentHandler) GetPaymentsByLoanID(ctx *fiber.Ctx) error {

// 	loanId, err := strconv.ParseInt(ctx.Params("loan_id"), 10, 64)
// 	if err != nil {
// 		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Loan ID"})
// 	}

// 	_status, err := strconv.ParseInt(ctx.Params("status"), 10, 8)
// 	if err != nil {
// 		return err
// 	}
// 	status := entity.PaymentStatus(_status)

// 	payments, err := h.paymentUsecase.GetPaymentsByLoanID(ctx.Context(), loanId, &status)
// 	if err != nil {
// 		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	if payments == nil {
// 		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"data": []entity.Payment{}})
// 	}

// 	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": payments})
// }
