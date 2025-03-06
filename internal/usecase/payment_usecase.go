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

func (u *PaymentUsecase) CreatePaymentsWithTx(tx *sql.Tx, payment []*entity.Payment) error {
	return u.paymentRepo.CreatePaymentsWithTx(tx, payment)
}
func (u *PaymentUsecase) GetPaymentByID(ctx context.Context, id int64) (*entity.Payment, error) {
	return u.paymentRepo.GetPaymentByID(ctx, id)
}
func (u *PaymentUsecase) GetAllPayments(ctx context.Context, status *entity.PaymentStatus) ([]*entity.Payment, error) {
	return u.paymentRepo.GetAllPayments(ctx, status)
}
func (u *PaymentUsecase) GetPaymentsByLoanID(ctx context.Context, loanId int64, status *entity.PaymentStatus, dueBefore *time.Time) ([]*entity.Payment, error) {
	return u.paymentRepo.GetPaymentsByLoanID(ctx, loanId, status, dueBefore)
}

func (u *PaymentUsecase) CreatePaymentsInTx(tx *sql.Tx, payments []entity.CreatePaymentPayload) error {
	if len(payments) == 0 {
		return errors.New("no payments to create")
	}

	paymentEntities := make([]*entity.Payment, len(payments))
	for i, payload := range payments {
		if err := u.validatePaymentPayload(payload); err != nil {
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

	return u.paymentRepo.CreatePaymentsWithTx(tx, paymentEntities)
}

func (u *PaymentUsecase) PayPayment(tx *sql.Tx, paymentID int64, transactionID int64, paidAt time.Time) error {
	return u.paymentRepo.PayPayment(tx, paymentID, transactionID, paidAt)
}

func (u *PaymentUsecase) validatePaymentPayload(req entity.CreatePaymentPayload) error {
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
