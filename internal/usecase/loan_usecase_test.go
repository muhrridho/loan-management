package usecase

import (
	"context"
	"loan-management/internal/entity"
	internalMock "loan-management/internal/mock"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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
var mockUserUsecase = new(internalMock.MockUserUsecase)

func setupMocks() (*internalMock.MockLoanRepository, *internalMock.MockUserUsecase, *LoanUsecase) {
	mockRepo := new(internalMock.MockLoanRepository)
	mockUserUsecase := new(internalMock.MockUserUsecase)

	// Set default mock expectations
	mockUserUsecase.On("GetByID", mock.Anything, mock.Anything).Return(nil, nil)
	mockUserUsecase.On("IsUserDelinquent", mock.Anything, mock.Anything).Return(false, nil)

	mockUsecase := NewLoanUsecase(mockRepo, mockUserUsecase)

	return mockRepo, mockUserUsecase, mockUsecase
}

func TestCreateLoan(t *testing.T) {

	t.Run("Success Create Loan", func(t *testing.T) {
		mockRepo, _, mockUsecase := setupMocks()

		mockRepo.On("GetByUserID", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
		mockRepo.On("Create", mock.Anything, mockLoan).Return(nil)

		err := mockUsecase.Create(context.Background(), mockLoan)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)

	})

	t.Run("Failed Create Loan - Billing Start At Invalid", func(t *testing.T) {
		mockRepo, _, mockUsecase := setupMocks()
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
		mockRepo, _, mockUsecase := setupMocks()

		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(mockLoan, nil)

		loan, err := mockUsecase.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, mockLoan, loan)
		mockRepo.AssertExpectations(t)
	})
}

func TestLoanGetAll(t *testing.T) {

	t.Run("Success Get All Loans", func(t *testing.T) {
		mockRepo, _, mockUsecase := setupMocks()
		expectedLoans := []*entity.Loan{mockLoan}

		mockRepo.On("GetAll", mock.Anything).Return(expectedLoans, nil)

		loans, err := mockUsecase.GetAll(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expectedLoans, loans)
		mockRepo.AssertExpectations(t)
	})
}
