package usecase

import (
	"context"
	"database/sql"
	"errors"
	"loan-management/internal/entity"
	"loan-management/internal/repository"
	"os"
	"strconv"
	"time"
)

var (
	ErrInvalidBillingStartDate = errors.New("billing start date cannot be in the past")
	ErrStillHasActiveLoan      = errors.New("Can't create loan because you still have an active loans")
)

type LoanUsecaseInterface interface {
	GetAllLoans(ctx context.Context) ([]*entity.Loan, error)
	GetLoanByID(ctx context.Context, id int64, status *entity.LoanStatus) (*entity.Loan, error)
	GetLoansByUserID(ctx context.Context, userID int64, status entity.LoanStatus) ([]*entity.Loan, error)
	CheckCreateLoanEligibility(ctx context.Context, loan *entity.Loan) error
	CreateLoanWithPayments(ctx context.Context, loan *entity.Loan) error
	GetLoanDuePayments(ctx context.Context, loan *entity.Loan) ([]*entity.Payment, error)
	UpdateLoanOutstanding(tx *sql.Tx, outstanding float64, loanID int64) error
}

type LoanUsecase struct {
	loanRepo       repository.LoanRepository
	userUsecase    UserUsecaseInterface
	paymentUsecase PaymentUsecaseInterface
}

func NewLoanUsecase(loanRepo repository.LoanRepository, userUsecase UserUsecaseInterface, paymentUsecase PaymentUsecaseInterface) *LoanUsecase {
	return &LoanUsecase{
		loanRepo:       loanRepo,
		userUsecase:    userUsecase,
		paymentUsecase: paymentUsecase,
	}
}

func (u *LoanUsecase) GetAllLoans(ctx context.Context) ([]*entity.Loan, error) {
	return u.loanRepo.GetAllLoans(ctx)
}

func (u *LoanUsecase) GetLoanByID(ctx context.Context, id int64, status *entity.LoanStatus) (*entity.Loan, error) {
	return u.loanRepo.GetLoanByID(ctx, id, status)
}

func (u *LoanUsecase) GetLoansByUserID(ctx context.Context, userID int64, status entity.LoanStatus) ([]*entity.Loan, error) {
	return u.loanRepo.GetLoansByUserID(ctx, userID, &status)
}

func (u *LoanUsecase) CheckCreateLoanEligibility(ctx context.Context, loan *entity.Loan) error {

	// Strict user to only have 1 active loan at one time
	// activeStatus := entity.LoanStatusActive
	// if loans, _ := u.GetLoansByUserID(ctx, loan.UserID, activeStatus); loans != nil {
	// 	return ErrStillHasActiveLoan
	// }

	// check if user is delinquent
	isUserDelinquent, err := u.userUsecase.IsUserDelinquent(ctx, loan.UserID)
	if err != nil {
		return err
	}
	if isUserDelinquent {
		return errors.New("Can't create loan due to user is delinquent")
	}

	return nil

}

func (u *LoanUsecase) CreateLoanWithPayments(ctx context.Context, loan *entity.Loan) error {
	if err := u.validateBillingStartDate(loan.BillingStartDate); err != nil {
		return err
	}

	if _, err := u.userUsecase.GetUserByID(ctx, loan.UserID); err != nil {
		return err
	}

	if err := u.CheckCreateLoanEligibility(ctx, loan); err != nil {
		return err
	}

	if loan.InterestType == entity.InterestTypeFlatAnnual {
		loan.Outstanding = u.calculateTotalOutstanding(loan)
	} else {
		// TODO: Implement Reduce Annual Type
	}

	tx, err := u.loanRepo.BeginTx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	loan, err = u.loanRepo.CreateLoan(tx, loan)
	if err != nil {
		return err
	}

	// calculate payment details
	totalInterest := u.calculateInterest(loan)
	amountPerInstallment := loan.Amount / float64(loan.Tenure)
	interestPerInstallment := totalInterest / float64(loan.Tenure)

	// generate payments
	paymentsPayload := make([]entity.CreatePaymentPayload, loan.Tenure)
	for i := 0; i < loan.Tenure; i++ {

		dueDate := loan.BillingStartDate.AddDate(0, 0, ((i + 1) * 7))

		paymentsPayload[i] = entity.CreatePaymentPayload{
			LoanID:      loan.ID,
			DueDate:     dueDate,
			PaymentNo:   int32(i + 1),
			Amount:      amountPerInstallment,
			Interest:    interestPerInstallment,
			TotalAmount: amountPerInstallment + interestPerInstallment,
		}
	}

	err = u.paymentUsecase.CreatePayment(tx, paymentsPayload)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (u *LoanUsecase) GetLoanDuePayments(ctx context.Context, loan *entity.Loan) ([]*entity.Payment, error) {

	var dueBefore time.Time
	if loan.TenureType == entity.TenureTypeWeekly {
		// added 7 days to include next due payments
		dueBefore = time.Now().AddDate(0, 0, 7)
	} else {
		// TODO: implement monthly calculation
	}

	paymentStatusActive := entity.PaymentStatusActive
	payments, err := u.paymentUsecase.GetPaymentsByLoanID(ctx, loan.ID, &paymentStatusActive, &dueBefore)

	if err != nil {
		return nil, err
	}

	return payments, nil
}

func (u *LoanUsecase) UpdateLoanOutstanding(tx *sql.Tx, outstanding float64, loanID int64) error {
	return u.loanRepo.UpdateLoanOutstanding(tx, outstanding, loanID)
}

func (u *LoanUsecase) validateBillingStartDate(billingStartDate time.Time) error {
	// for testing purpose: enable loan creating with start billing date that already in the past
	allowPastDate, err := strconv.ParseBool(os.Getenv("ALLOW_CREATE_LOAN_PAST_DATE"))
	if err != nil {
		allowPastDate = false
	}
	if !allowPastDate && billingStartDate.Before(time.Now().Truncate(24*time.Hour)) {
		return ErrInvalidBillingStartDate
	}
	return nil
}

func (u *LoanUsecase) calculateTotalOutstanding(loan *entity.Loan) float64 {
	totalInterest := u.calculateInterest(loan)

	totalOutstanding := loan.Amount + totalInterest

	return totalOutstanding
}

func (u *LoanUsecase) calculateInterest(loan *entity.Loan) float64 {
	tenureInYears := float64(loan.Tenure) / 52
	// if loan.TenureType == entity.TenureTypeMonthly {
	// 	tenureInYears = float64(loan.Tenure) / 12
	// }
	return loan.Amount * (loan.Interest / 100) * tenureInYears
}
