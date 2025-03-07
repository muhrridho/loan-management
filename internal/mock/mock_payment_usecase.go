package mock

import (
	"context"
	"database/sql"
	"loan-management/internal/entity"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockPaymentUsecase struct {
	mock.Mock
}

func (m *MockPaymentUsecase) CreatePayment(tx *sql.Tx, payments []entity.CreatePaymentPayload) error {
	args := m.Called(tx, payments)
	return args.Error(0)
}

func (m *MockPaymentUsecase) GetPaymentByID(ctx context.Context, id int64) (*entity.Payment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.Payment), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPaymentUsecase) GetAllPayments(ctx context.Context, status *entity.PaymentStatus) ([]*entity.Payment, error) {
	args := m.Called(ctx, status)
	if args.Get(0) != nil {
		return args.Get(0).([]*entity.Payment), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPaymentUsecase) GetPaymentsByLoanID(ctx context.Context, loanId int64, status *entity.PaymentStatus, dueBefore *time.Time) ([]*entity.Payment, error) {
	args := m.Called(ctx, loanId, status, dueBefore)
	if args.Get(0) != nil {
		return args.Get(0).([]*entity.Payment), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPaymentUsecase) CreatePaymentsInTx(tx *sql.Tx, payments []entity.CreatePaymentPayload) error {
	args := m.Called(tx, payments)
	return args.Error(0)
}

func (m *MockPaymentUsecase) PayPayment(tx *sql.Tx, paymentID int64, transactionID int64, paidAt time.Time) error {
	args := m.Called(tx, paymentID, transactionID, paidAt)
	return args.Error(0)
}
