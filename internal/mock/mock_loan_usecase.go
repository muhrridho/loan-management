package mock

import (
	"context"
	"database/sql"
	"loan-management/internal/entity"

	"github.com/stretchr/testify/mock"
)

type MockLoanUsecase struct {
	mock.Mock
}

func (m *MockLoanUsecase) GetAllLoans(ctx context.Context) ([]*entity.Loan, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]*entity.Loan), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockLoanUsecase) GetLoanByID(ctx context.Context, id int64, status *entity.LoanStatus) (*entity.Loan, error) {
	args := m.Called(ctx, id, status)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.Loan), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockLoanUsecase) GetLoansByUserID(ctx context.Context, userID int64, status entity.LoanStatus) ([]*entity.Loan, error) {
	args := m.Called(ctx, userID, status)
	if args.Get(0) != nil {
		return args.Get(0).([]*entity.Loan), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockLoanUsecase) CheckCreateLoanEligibility(ctx context.Context, loan *entity.Loan) error {
	args := m.Called(ctx, loan)
	return args.Error(0)
}

func (m *MockLoanUsecase) CreateLoanWithPayments(ctx context.Context, loan *entity.Loan) error {
	args := m.Called(ctx, loan)
	return args.Error(0)
}

func (m *MockLoanUsecase) GetLoanDuePayments(ctx context.Context, loan *entity.Loan) ([]*entity.Payment, error) {
	args := m.Called(ctx, loan)
	if args.Get(0) != nil {
		return args.Get(0).([]*entity.Payment), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockLoanUsecase) UpdateLoanOutstanding(tx *sql.Tx, outstanding float64, loanID int64) error {
	args := m.Called(tx, outstanding, loanID)
	return args.Error(0)
}
