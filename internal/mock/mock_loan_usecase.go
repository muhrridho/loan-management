package mock

import (
	"context"
	"loan-management/internal/entity"

	"github.com/stretchr/testify/mock"
)

type MockLoanUsecase struct {
	mock.Mock
}

func (m *MockLoanUsecase) CreateLoan(ctx context.Context, loan *entity.Loan) error {
	args := m.Called(ctx, loan)
	return args.Error(0)
}

func (m *MockLoanUsecase) GetAllLoans(ctx context.Context) ([]*entity.Loan, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Loan), args.Error(1)
}

func (m *MockLoanUsecase) GetLoanByID(ctx context.Context, id int64) (*entity.Loan, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Loan), args.Error(1)
}

func (m *MockLoanUsecase) GetLoansByUserID(ctx context.Context, userID int64, status entity.LoanStatus) ([]*entity.Loan, error) {
	args := m.Called(ctx, userID, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Loan), args.Error(1)
}

func (m *MockLoanUsecase) CheckCreateLoanEligibility(ctx context.Context, loan *entity.Loan) error {
	args := m.Called(ctx, loan)
	return args.Error(0)
}
