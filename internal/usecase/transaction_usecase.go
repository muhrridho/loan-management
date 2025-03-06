package usecase

import (
	"context"
	"errors"
	"fmt"
	"loan-management/internal/entity"
	"loan-management/internal/repository"
	"time"
)

type TransactionUsecase struct {
	transactionRepository repository.TransactionRepository
	loanUsecase           LoanUsecase
	paymentUsecase        PaymentUsecase
}

func NewTransactionUsecase(transactionRepository repository.TransactionRepository, loanUsecase *LoanUsecase, paymentUsecase *PaymentUsecase) *TransactionUsecase {
	return &TransactionUsecase{
		transactionRepository: transactionRepository,
		loanUsecase:           *loanUsecase,
		paymentUsecase:        *paymentUsecase,
	}
}

func (u *TransactionUsecase) InquiryTransaction(ctx context.Context, loanID int64) (*entity.TransactionInquiry, error) {
	loanStatusActive := entity.LoanStatusActive
	loan, err := u.loanUsecase.GetLoanByID(ctx, loanID, &loanStatusActive)
	if err != nil {
		return nil, err
	}

	if loan == nil {
		return nil, nil
	}

	duePayments, err := u.loanUsecase.GetLoanDuePayments(ctx, loan)

	if err != nil {
		return nil, err
	}

	if len(duePayments) <= 0 {
		return nil, nil
	}

	var amountDue float64
	for _, payment := range duePayments {
		amountDue += payment.TotalAmount
	}

	// assume the payments is in asc order
	latestDueDate := duePayments[len(duePayments)-1].DueDate

	inquiryResult := &entity.TransactionInquiry{
		LoanID:     loan.ID,
		AmountDue:  amountDue,
		DueDate:    latestDueDate,
		LoanDetail: loan,
		Bills:      duePayments,
	}

	return inquiryResult, nil
}

func (u *TransactionUsecase) CreateTransaction(ctx context.Context, trxPayload *entity.CreateTransactionPayload) (*entity.Transaction, error) {

	// get active loan based on payload LoanID
	loanStatusActive := entity.LoanStatusActive
	loan, err := u.loanUsecase.GetLoanByID(ctx, trxPayload.LoanID, &loanStatusActive)

	if err != nil {
		return nil, err
	}

	if loan == nil {
		return nil, nil
	}

	// get all due payments that will be paid in this trx
	duePayments, err := u.loanUsecase.GetLoanDuePayments(ctx, loan)
	fmt.Println(duePayments)

	if err != nil {
		return nil, err
	}

	if len(duePayments) <= 0 {
		return nil, nil
	}

	var amountDue float64
	for _, payment := range duePayments {
		amountDue += payment.TotalAmount
	}

	// validate amount
	if amountDue != trxPayload.Amount {
		return nil, errors.New("The amount is different with the due amount")
	}

	/**
	 * Begin the DB trx; steps:
	 * 1. Create transaction
	 * 2. Pay all payments (set status = paid, etc)
	 * 3. Update Loan (outstanding, status, etc)
	 */
	tx, err := u.transactionRepository.BeginTx()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Create transaction step
	timeNow := time.Now()
	trxStatusPaid := entity.TransactionStatusPaid

	// Assume no waiting for payment, so trx will be set directly as paid
	trx := &entity.Transaction{
		TotalAmount: amountDue,
		Penalty:     0,
		Status:      trxStatusPaid,
		PaidAt:      &timeNow,
		CreatedAt:   timeNow,
	}

	trxID, err := u.transactionRepository.CreateTransaction(tx, trx)

	if err != nil {
		return nil, err
	}

	if trxID == 0 {
		return nil, errors.New("Something went wrong")
	}

	trx.ID = trxID

	// Pay all payments step
	for _, payment := range duePayments {
		err := u.paymentUsecase.PayPayment(tx, payment.ID, trxID, timeNow)
		if err != nil {
			return nil, err
		}
	}

	// Update loan step
	outstanding := loan.Outstanding - amountDue
	if err := u.loanUsecase.UpdateLoanOutstanding(tx, outstanding, loan.ID); err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return trx, nil
}
