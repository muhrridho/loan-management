package mock

import (
	"context"
	"database/sql"
	"loan-management/internal/entity"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockPaymentRepository struct {
	mock.Mock
}

func (m *MockPaymentRepository) CreatePayment(tx *sql.Tx, payments []*entity.Payment) error {
	args := m.Called(tx, payments)
	return args.Error(0)
}

func (m *MockPaymentRepository) GetPaymentByID(ctx context.Context, id int64) (*entity.Payment, error) {
	args := m.Called(ctx, id)
	if payment, ok := args.Get(0).(*entity.Payment); ok {
		return payment, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPaymentRepository) GetAllPayments(ctx context.Context, status *entity.PaymentStatus) ([]*entity.Payment, error) {
	args := m.Called(ctx, status)
	if payments, ok := args.Get(0).([]*entity.Payment); ok {
		return payments, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPaymentRepository) GetPaymentsByLoanID(ctx context.Context, loanId int64, status *entity.PaymentStatus, dueBefore *time.Time) ([]*entity.Payment, error) {
	args := m.Called(ctx, loanId, status, dueBefore)
	if payments, ok := args.Get(0).([]*entity.Payment); ok {
		return payments, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPaymentRepository) PayPayment(tx *sql.Tx, paymentId int64, transactionId int64, paidAt time.Time) error {
	args := m.Called(tx, paymentId, transactionId, paidAt)
	return args.Error(0)
}
