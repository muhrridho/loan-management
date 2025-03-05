package mock

import (
	"context"
	"loan-management/internal/entity"

	"github.com/stretchr/testify/mock"
)

type MockLoanRepository struct {
	mock.Mock
}

func (m *MockLoanRepository) CreateLoan(ctx context.Context, loan *entity.Loan) error {
	args := m.Called(ctx, loan)
	return args.Error(0)
}

func (m *MockLoanRepository) GetAllLoans(ctx context.Context) ([]*entity.Loan, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entity.Loan), args.Error(1)
}

func (m *MockLoanRepository) GetLoanByID(ctx context.Context, id int64) (*entity.Loan, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.Loan), args.Error(1)
}

func (m *MockLoanRepository) GetLoansByUserID(ctx context.Context, userID int64, status *entity.LoanStatus) ([]*entity.Loan, error) {
	args := m.Called(ctx, userID, status)
	if args.Get(0) != nil {
		return args.Get(0).([]*entity.Loan), args.Error(1)
	}
	return nil, args.Error(1)
}
