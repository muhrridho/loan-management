package usecase

import (
	"context"
	"database/sql"
	"loan-management/internal/entity"
	internalMock "loan-management/internal/mock"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var MockPayment = &entity.Payment{
	LoanID:        1,
	TransactionID: nil,
	PaymentNo:     1,
	DueDate:       time.Time{},
	Amount:        float64(1000000),
	Interest:      float64(100000),
	TotalAmount:   float64(1100000),
	Status:        entity.PaymentStatusActive,
	PaidAt:        &time.Time{},
	CreatedAt:     time.Time{},
}

func TestGetPaymentByID(t *testing.T) {
	t.Run("Success Get Payment By ID", func(t *testing.T) {
		mockRepo := new(internalMock.MockPaymentRepository)
		mockUsecase := NewPaymentUsecase(mockRepo)
		mockRepo.On("GetPaymentByID", mock.Anything, mock.Anything).Return(MockPayment, nil)

		payment, err := mockUsecase.GetPaymentByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, payment, MockPayment)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetAllPayments(t *testing.T) {
	t.Run("Success GetAllPayments", func(t *testing.T) {
		mockRepo := new(internalMock.MockPaymentRepository)
		mockUsecase := NewPaymentUsecase(mockRepo)

		mockPayments := []*entity.Payment{MockPayment}
		mockRepo.On("GetAllPayments", mock.Anything, mock.Anything).Return(mockPayments, nil)
		paymentStatusActive := entity.PaymentStatusActive
		payments, err := mockUsecase.GetAllPayments(context.Background(), &paymentStatusActive)

		assert.NoError(t, err)
		assert.Equal(t, payments, mockPayments)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetPaymentsByLoanID(t *testing.T) {
	t.Run("Success GetPaymentsByLoanID", func(t *testing.T) {
		mockRepo := new(internalMock.MockPaymentRepository)
		mockUsecase := NewPaymentUsecase(mockRepo)

		mockPayments := []*entity.Payment{MockPayment}
		mockRepo.On("GetPaymentsByLoanID", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockPayments, nil)
		paymentStatusActive := entity.PaymentStatusActive
		payments, err := mockUsecase.GetPaymentsByLoanID(context.Background(), int64(1), &paymentStatusActive, &time.Time{})

		assert.NoError(t, err)
		assert.Equal(t, payments, mockPayments)
		mockRepo.AssertExpectations(t)
	})
}

func TestCreatePayment(t *testing.T) {
	t.Run("Success CreatePayment", func(t *testing.T) {
		mockRepo := new(internalMock.MockPaymentRepository)
		mockUsecase := NewPaymentUsecase(mockRepo)

		mockRepo.On("CreatePayment", mock.Anything, mock.Anything).Return(nil)
		mockCreatePaymentPayload := []entity.CreatePaymentPayload{{
			LoanID:      1,
			DueDate:     time.Now(),
			PaymentNo:   1,
			Amount:      1000000,
			Interest:    100000,
			TotalAmount: 1100000,
		}}
		err := mockUsecase.CreatePayment(&sql.Tx{}, mockCreatePaymentPayload)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestPayPayment(t *testing.T) {
	t.Run("Success PayPayment", func(t *testing.T) {
		mockRepo := new(internalMock.MockPaymentRepository)
		mockUsecase := NewPaymentUsecase(mockRepo)

		mockRepo.On("PayPayment", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
		err := mockUsecase.PayPayment(&sql.Tx{}, int64(1), int64(1), time.Time{})

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// func Test(t *testing.T) {
// 	t.Run("Success ", func(t *testing.T) {
// 		mockRepo := new(internalMock.MockPaymentRepository)
// 		mockUsecase := NewPaymentUsecase(mockRepo)

// 		mockRepo.On("GetPaymentByID", mock.Anything, mock.Anything).Return(MockPayment, nil)
// 		payment, err := mockUsecase.GetPaymentByID(context.Background(), 1)

// 		assert.NoError(t, err)
// 		assert.Equal(t, payment, MockPayment)
// 		mockRepo.AssertExpectations(t)
// 	})
// }
