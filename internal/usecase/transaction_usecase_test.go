package usecase

import (
	"context"
	"errors"
	"loan-management/internal/entity"
	internalMock "loan-management/internal/mock"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTransactionUsecase struct {
	mock.Mock
}

var MockTransaction = &entity.Transaction{
	ID:          1,
	TotalAmount: 5500000,
	Penalty:     0,
	Status:      entity.TransactionStatusActive,
	PaidAt:      &time.Time{},
	CreatedAt:   time.Time{},
}

func setupTransactionMocks() (*TransactionUsecase, *internalMock.MockTransactionRepository, *internalMock.MockLoanUsecase, *internalMock.MockPaymentUsecase) {
	mockRepo := new(internalMock.MockTransactionRepository)
	mockLoanUsecase := new(internalMock.MockLoanUsecase)
	mockPaymentUsecase := new(internalMock.MockPaymentUsecase)

	mockUsecase := NewTransactionUsecase(mockRepo, mockLoanUsecase, mockPaymentUsecase)

	return mockUsecase, mockRepo, mockLoanUsecase, mockPaymentUsecase
}

func TestInquiryTransaction(t *testing.T) {
	t.Run("Success InquiryTransaction", func(t *testing.T) {
		mockUsecase, mockRepo, mockLoanUsecase, _ := setupTransactionMocks()

		mockPayments := []*entity.Payment{MockPayment}
		mockTransactionInquiry := entity.TransactionInquiry{
			LoanID:     int64(0),
			AmountDue:  float64(1100000),
			DueDate:    time.Time{},
			LoanDetail: MockLoan,
			Bills:      mockPayments,
		}
		mockLoanUsecase.On("GetLoanByID", mock.Anything, mock.Anything, mock.Anything).Return(mockTransactionInquiry.LoanDetail, nil)
		mockLoanUsecase.On("GetLoanDuePayments", mock.Anything, mock.Anything).Return(mockPayments, nil)
		inquiryResult, err := mockUsecase.InquiryTransaction(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, *inquiryResult, mockTransactionInquiry)
		mockRepo.AssertExpectations(t)
	})
}

func TestCreateTransaction(t *testing.T) {

	createTrxPayload := entity.CreateTransactionPayload{
		LoanID: MockTransaction.ID,
		Amount: MockPayment.TotalAmount,
	}

	t.Run("Success CreateTransaction", func(t *testing.T) {
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

		defer func() { now = time.Now }()

		mockUsecase, mockRepo, mockLoanUsecase, mockPaymentUsecase := setupTransactionMocks()

		mockPayments := []*entity.Payment{MockPayment}
		mockLoanUsecase.On("GetLoanByID", mock.Anything, mock.Anything, mock.Anything).Return(MockLoan, nil)
		mockLoanUsecase.On("GetLoanDuePayments", mock.Anything, mock.Anything).Return(mockPayments, nil)
		mockLoanUsecase.On("UpdateLoanOutstanding", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		mockRepo.On("BeginTx").Return(mockTx, nil)
		mockRepo.On("CreateTransaction", mock.Anything, mock.Anything).Return(int64(1), nil)

		mockPaymentUsecase.On("PayPayment", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		trx, err := mockUsecase.CreateTransaction(context.Background(), &createTrxPayload)
		mockPaidTransaction := MockTransaction
		mockPaidTransaction.Status = entity.TransactionStatusPaid
		mockPaidTransaction.CreatedAt = mockTime
		mockPaidTransaction.PaidAt = &mockTime
		mockPaidTransaction.TotalAmount = MockPayment.TotalAmount

		assert.NoError(t, err)
		assert.Equal(t, trx, MockTransaction)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed CreateTransaction - Loan Not Found", func(t *testing.T) {
		mockUsecase, mockRepo, mockLoanUsecase, _ := setupTransactionMocks()
		mockLoanUsecase.On("GetLoanByID", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

		trx, err := mockUsecase.CreateTransaction(context.Background(), &createTrxPayload)

		assert.NoError(t, err)
		assert.Equal(t, trx, (*entity.Transaction)(nil))
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed CreateTransaction - No Due Payment", func(t *testing.T) {
		mockUsecase, mockRepo, mockLoanUsecase, _ := setupTransactionMocks()
		mockPayments := []*entity.Payment{MockPayment}
		mockLoanUsecase.On("GetLoanByID", mock.Anything, mock.Anything, mock.Anything).Return(MockLoan, nil)
		mockLoanUsecase.On("GetLoanDuePayments", mock.Anything, mock.Anything).Return(mockPayments, nil)
		mockLoanUsecase.On("UpdateLoanOutstanding", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		trx, err := mockUsecase.CreateTransaction(context.Background(), &createTrxPayload)

		assert.NoError(t, err)
		assert.Equal(t, trx, (*entity.Transaction)(nil))
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed CreateTransaction - Total Amount Not Match", func(t *testing.T) {
		mockUsecase, mockRepo, mockLoanUsecase, _ := setupTransactionMocks()

		mockPayments := []*entity.Payment{MockPayment}
		mockLoanUsecase.On("GetLoanByID", mock.Anything, mock.Anything, mock.Anything).Return(MockLoan, nil)
		mockLoanUsecase.On("GetLoanDuePayments", mock.Anything, mock.Anything).Return(mockPayments, nil)
		mockLoanUsecase.On("UpdateLoanOutstanding", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		customCreateTrxPayload := createTrxPayload
		customCreateTrxPayload.Amount = 500

		trx, err := mockUsecase.CreateTransaction(context.Background(), &customCreateTrxPayload)

		assert.Error(t, err)
		assert.Equal(t, err, errors.New("The amount is different with the due amount"))
		assert.Equal(t, trx, (*entity.Transaction)(nil))
		mockRepo.AssertExpectations(t)
	})
}
