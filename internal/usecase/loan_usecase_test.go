package usecase

import (
	"context"
	"database/sql"
	"errors"
	"loan-management/internal/entity"
	internalMock "loan-management/internal/mock"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var MockLoan = &entity.Loan{
	UserID:           1,
	Interest:         10,
	InterestType:     entity.InterestTypeFlatAnnual,
	Tenure:           1,
	TenureType:       entity.TenureTypeWeekly,
	Status:           entity.LoanStatusActive,
	Amount:           1000000,
	CreatedAt:        time.Now(),
	BillingStartDate: time.Now(),
}

func setupMocks() (*internalMock.MockLoanRepository, *internalMock.MockUserUsecase, *internalMock.MockPaymentUsecase, *LoanUsecase) {
	mockRepo := new(internalMock.MockLoanRepository)
	mockUserUsecase := new(internalMock.MockUserUsecase)
	mockPaymentUsecase := new(internalMock.MockPaymentUsecase)

	mockUsecase := NewLoanUsecase(mockRepo, mockUserUsecase, mockPaymentUsecase)

	return mockRepo, mockUserUsecase, mockPaymentUsecase, mockUsecase
}

func TestGetAllLoans(t *testing.T) {

	t.Run("Success GetAllLoans", func(t *testing.T) {
		mockRepo, _, _, mockUsecase := setupMocks()
		expectedLoans := []*entity.Loan{MockLoan}

		mockRepo.On("GetAllLoans", mock.Anything).Return(expectedLoans, nil)

		loans, err := mockUsecase.GetAllLoans(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expectedLoans, loans)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetLoanByID(t *testing.T) {
	t.Run("Success GetLoanByID", func(t *testing.T) {
		mockRepo, _, _, mockUsecase := setupMocks()

		loanStatusActive := entity.LoanStatusActive
		mockRepo.On("GetLoanByID", mock.Anything, mock.Anything, mock.Anything).Return(MockLoan, nil)

		loan, err := mockUsecase.GetLoanByID(context.Background(), 1, &loanStatusActive)

		assert.NoError(t, err)
		assert.Equal(t, MockLoan, loan)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetLoansByUserID(t *testing.T) {
	t.Run("Success GetLoanByUserID", func(t *testing.T) {
		mockRepo, _, _, mockUsecase := setupMocks()

		loanStatusActive := entity.LoanStatusActive
		mockLoans := []*entity.Loan{MockLoan}
		mockRepo.On("GetLoansByUserID", mock.Anything, mock.Anything, mock.Anything).Return(mockLoans, nil)

		loans, err := mockUsecase.GetLoansByUserID(context.Background(), 1, loanStatusActive)

		assert.NoError(t, err)
		assert.Equal(t, mockLoans, loans)
		mockRepo.AssertExpectations(t)
	})
}

func TestCheckCreateLoanEligibility(t *testing.T) {
	t.Run("Success CheckCreateLoanEligibility ", func(t *testing.T) {
		mockRepo, mockUserUsecase, _, mockUsecase := setupMocks()
		mockUserUsecase.On("IsUserDelinquent", mock.Anything, mock.Anything).Return(false, nil)

		err := mockUsecase.CheckCreateLoanEligibility(context.Background(), MockLoan)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)

	})
}

func TestCreateLoanWithPayments(t *testing.T) {

	t.Run("Success CreateLoan", func(t *testing.T) {
		// setup dbMock to get mock of sql.tx
		db, dbMock, _ := sqlmock.New()
		defer db.Close()
		dbMock.ExpectBegin()
		mockTx, _ := db.Begin()
		dbMock.ExpectCommit()

		mockTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		now = func() time.Time {
			return mockTime
		}

		mockRepo, mockUserUsecase, mockPaymentUsecase, mockUsecase := setupMocks()

		mockRepo.On("CreateLoan", mock.Anything, mock.Anything).Return(MockLoan, nil)
		mockRepo.On("BeginTx").Return(mockTx, nil)
		mockUserUsecase.On("IsUserDelinquent", mock.Anything, mock.Anything).Return(false, nil)
		mockUserUsecase.On("GetUserByID", mock.Anything, mock.Anything).Return(MockUser, nil)

		expectedInterest := float64((MockLoan.Amount * (MockLoan.Interest / 100)) / 52)
		customMockPaymentPayload := entity.CreatePaymentPayload{
			LoanID:      MockLoan.ID,
			DueDate:     MockLoan.BillingStartDate.AddDate(0, 0, 7),
			PaymentNo:   int32(1),
			Amount:      MockLoan.Amount,
			Interest:    expectedInterest,
			TotalAmount: MockLoan.Amount/float64(MockLoan.Tenure) + expectedInterest,
		}
		expectedPaymentPayload := []*entity.CreatePaymentPayload{&customMockPaymentPayload}
		mockPaymentUsecase.On("CreatePayment", mock.Anything, mock.MatchedBy(func(payloads []entity.CreatePaymentPayload) bool {
			t.Logf("Payloads: %+v", expectedPaymentPayload)

			// Convert payloads to slice of pointers
			payloadPtrs := make([]*entity.CreatePaymentPayload, len(payloads))
			for i := range payloads {
				payloadPtrs[i] = &payloads[i]
			}

			// ensuring correct payloads when creating payments
			return assert.EqualValues(t, expectedPaymentPayload, payloadPtrs)
		})).Return(nil)
		err := mockUsecase.CreateLoanWithPayments(context.Background(), MockLoan)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed CreateLoan - User Not Found", func(t *testing.T) {
		mockRepo, mockUserUsecase, _, mockUsecase := setupMocks()
		mockUserUsecase.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, errors.New(""))
		err := mockUsecase.CreateLoanWithPayments(context.Background(), MockLoan)
		assert.Error(t, err)
		assert.Equal(t, err, errors.New(""))
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed CreateLoan - User Not Eligible", func(t *testing.T) {
		mockRepo, mockUserUsecase, _, mockUsecase := setupMocks()
		mockUserUsecase.On("GetUserByID", mock.Anything, mock.Anything).Return(MockUser, nil)
		mockUserUsecase.On("IsUserDelinquent", mock.Anything, mock.Anything).Return(true, errors.New(""))
		err := mockUsecase.CreateLoanWithPayments(context.Background(), MockLoan)
		assert.Error(t, err)
		assert.Equal(t, err, errors.New(""))
		mockRepo.AssertExpectations(t)
	})

}

func TestGetLoanDuePayments(t *testing.T) {
	t.Run("Success GetLoanDuePayments", func(t *testing.T) {
		mockRepo, _, mockPaymentUsecase, mockUsecase := setupMocks()

		mockPayments := []*entity.Payment{MockPayment}
		mockPaymentUsecase.On("GetPaymentsByLoanID", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockPayments, nil)

		payments, err := mockUsecase.GetLoanDuePayments(context.Background(), MockLoan)
		assert.NoError(t, err)
		assert.Equal(t, mockPayments, payments)
		mockRepo.AssertExpectations(t)

	})
}

func TestUpdateLoanOutstanding(t *testing.T) {
	t.Run("Success UpdateLoanOutstanding", func(t *testing.T) {

		mockRepo, _, _, mockUsecase := setupMocks()
		mockRepo.On("UpdateLoanOutstanding", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		outstanding := float64(69)
		err := mockUsecase.UpdateLoanOutstanding(&sql.Tx{}, outstanding, 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)

	})
}

// func Test(t *testing.T) {
// 	t.Run("Success ", func(t *testing.T) {

// 	})
// }
