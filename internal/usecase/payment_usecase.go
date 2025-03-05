package usecase

import (
	"context"
	"loan-management/internal/entity"
	"loan-management/internal/repository"
)

type PaymentUsecase struct {
	paymentRepo repository.PaymentRepository
}

func NewPaymentUsecase(paymentRepo repository.PaymentRepository) *PaymentUsecase {
	return &PaymentUsecase{
		paymentRepo: paymentRepo,
	}
}

func (uc *PaymentUsecase) CreatePayment(ctx context.Context, payment *entity.Payment) error {
	return uc.paymentRepo.CreatePayment(ctx, payment)
}
func (uc *PaymentUsecase) GetPaymentByID(ctx context.Context, id int64) (*entity.Payment, error) {
	return uc.paymentRepo.GetPaymentByID(ctx, id)
}
func (uc *PaymentUsecase) GetAllPayments(ctx context.Context, status *entity.PaymentStatus) ([]*entity.Payment, error) {
	return uc.paymentRepo.GetAllPayments(ctx, status)
}
func (uc *PaymentUsecase) GetPaymentByLoanID(ctx context.Context, loanId int64, status *entity.PaymentStatus) ([]*entity.Payment, error) {
	return uc.paymentRepo.GetPaymentsByLoanID(ctx, loanId, status)
}
