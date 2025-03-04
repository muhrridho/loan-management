package usecase

import (
	"context"
	"loan-management/internal/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockLoanRepository struct {
	mock.Mock
}

func (m *MockLoanRepository) Create(ctx context.Context, loan *entity.Loan) error {
	args := m.Called(ctx, loan)
	return args.Error(0)
}

func (m *MockLoanRepository) GetAll(ctx context.Context) ([]*entity.Loan, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entity.Loan), args.Error(1)
}

func (m *MockLoanRepository) GetByID(ctx context.Context, id int64) (*entity.Loan, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.Loan), args.Error(1)
}

var mockLoan = &entity.Loan{
	UserID:           1,
	Interest:         10,
	InterestType:     entity.InterestTypeFlatAnnual,
	Tenure:           52,
	TenureType:       entity.TenureTypeWeekly,
	Status:           entity.LoanStatusActive,
	Amount:           5000000,
	CreatedAt:        time.Now(),
	BillingStartDate: time.Now(),
}

func TestCreateLoan(t *testing.T) {

	t.Run("Success Create Loan", func(t *testing.T) {
		mockRepo := new(MockLoanRepository)
		mockUsecase := NewLoanUsecase(mockRepo)
		mockRepo.On("Create", mock.Anything, mockLoan).Return(nil)

		err := mockUsecase.Create(context.Background(), mockLoan)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)

	})

	t.Run("Failed Create Loan - Billing Start At Invalid", func(t *testing.T) {
		mockRepo := new(MockLoanRepository)
		mockUsecase := NewLoanUsecase(mockRepo)
		customMockLoan := *mockLoan
		customMockLoan.BillingStartDate, _ = time.Parse("2006-01-02", "1945-08-01")

		mockRepo.On("Create", mock.Anything, customMockLoan).Return(nil)

		err := mockUsecase.Create(context.Background(), &customMockLoan)

		assert.Error(t, err)
		assert.Equal(t, "billing start date cannot be in the past", err.Error())

		mockRepo.AssertNotCalled(t, "Create")

	})
}

func TestLoanGetByID(t *testing.T) {

	t.Run("Success Get Loan By ID", func(t *testing.T) {
		mockRepo := new(MockLoanRepository)
		mockUsecase := NewLoanUsecase(mockRepo)

		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(mockLoan, nil)

		loan, err := mockUsecase.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, mockLoan, loan)
		mockRepo.AssertExpectations(t)
	})
}

func TestLoanGetAll(t *testing.T) {

	t.Run("Success Get All Loans", func(t *testing.T) {
		mockRepo := new(MockLoanRepository)
		mockUsecase := NewLoanUsecase(mockRepo)
		expectedLoans := []*entity.Loan{mockLoan}

		mockRepo.On("GetAll", mock.Anything).Return(expectedLoans, nil)

		loans, err := mockUsecase.GetAll(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expectedLoans, loans)
		mockRepo.AssertExpectations(t)
	})
}
