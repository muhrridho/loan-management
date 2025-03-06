package usecase

import (
	"context"
	"database/sql"
	"errors"
	"loan-management/internal/entity"
	"loan-management/internal/repository"
	"time"
)

type PaymentUsecase struct {
	paymentRepo repository.PaymentRepository
}

func NewPaymentUsecase(paymentRepo repository.PaymentRepository) *PaymentUsecase {
	return &PaymentUsecase{
		paymentRepo: paymentRepo,
	}
}

func (uc *PaymentUsecase) CreatePaymentsWithTx(tx *sql.Tx, payment []*entity.Payment) error {
	return uc.paymentRepo.CreatePaymentsWithTx(tx, payment)
}
func (uc *PaymentUsecase) GetPaymentByID(ctx context.Context, id int64) (*entity.Payment, error) {
	return uc.paymentRepo.GetPaymentByID(ctx, id)
}
func (uc *PaymentUsecase) GetAllPayments(ctx context.Context, status *entity.PaymentStatus) ([]*entity.Payment, error) {
	return uc.paymentRepo.GetAllPayments(ctx, status)
}
func (uc *PaymentUsecase) GetPaymentsByLoanID(ctx context.Context, loanId int64, status *entity.PaymentStatus, dueBefore *time.Time) ([]*entity.Payment, error) {
	return uc.paymentRepo.GetPaymentsByLoanID(ctx, loanId, status, dueBefore)
}

func (uc *PaymentUsecase) CreatePaymentsInTx(tx *sql.Tx, payments []entity.CreatePaymentPayload) error {
	if len(payments) == 0 {
		return errors.New("no payments to create")
	}

	paymentEntities := make([]*entity.Payment, len(payments))
	for i, payload := range payments {
		if err := uc.validatePaymentPayload(payload); err != nil {
			return err
		}

		paymentEntities[i] = &entity.Payment{
			LoanID:      payload.LoanID,
			DueDate:     payload.DueDate,
			PaymentNo:   payload.PaymentNo,
			Amount:      payload.Amount,
			Interest:    payload.Interest,
			TotalAmount: payload.TotalAmount,
			Status:      entity.PaymentStatusActive,
			PaidAt:      nil,
			CreatedAt:   time.Now(),
		}
	}

	return uc.paymentRepo.CreatePaymentsWithTx(tx, paymentEntities)
}

func (uc *PaymentUsecase) PayPayment(tx *sql.Tx, paymentID int64, transactionID int64, paidAt time.Time) error {
	return uc.paymentRepo.PayPayment(tx, paymentID, transactionID, paidAt)
}

func (uc *PaymentUsecase) validatePaymentPayload(req entity.CreatePaymentPayload) error {
	if req.LoanID <= 0 {
		return errors.New("invalid loan ID")
	}

	if req.PaymentNo <= 0 {
		return errors.New("invalid payment number")
	}

	if req.Amount <= 0 {
		return errors.New("payment amount must be positive")
	}

	if req.DueDate.IsZero() {
		return errors.New("due date is required")
	}

	return nil
}
