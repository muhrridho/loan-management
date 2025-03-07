package mock

import (
	"context"
	"database/sql"
	"loan-management/internal/entity"

	"github.com/stretchr/testify/mock"
)

type MockLoanRepository struct {
	mock.Mock
}

func (m *MockLoanRepository) BeginTx() (*sql.Tx, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).(*sql.Tx), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockLoanRepository) CreateLoan(tx *sql.Tx, loan *entity.Loan) (*entity.Loan, error) {
	args := m.Called(tx, loan)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.Loan), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockLoanRepository) GetLoanByID(ctx context.Context, id int64, status *entity.LoanStatus) (*entity.Loan, error) {
	args := m.Called(ctx, id, status)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.Loan), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockLoanRepository) GetAllLoans(ctx context.Context) ([]*entity.Loan, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]*entity.Loan), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockLoanRepository) GetLoansByUserID(ctx context.Context, userId int64, status *entity.LoanStatus) ([]*entity.Loan, error) {
	args := m.Called(ctx, userId, status)
	if args.Get(0) != nil {
		return args.Get(0).([]*entity.Loan), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockLoanRepository) UpdateLoanOutstanding(tx *sql.Tx, outstanding float64, loanID int64) error {
	args := m.Called(tx, outstanding, loanID)
	return args.Error(0)
}
