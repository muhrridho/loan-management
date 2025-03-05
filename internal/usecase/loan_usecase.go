package usecase

import (
	"context"
	"errors"
	"loan-management/internal/entity"
	"loan-management/internal/repository"
	"time"
)

var (
	ErrInvalidBillingStartDate = errors.New("billing start date cannot be in the past")
	ErrStillHasActiveLoan      = errors.New("Can't create loan because you still have an active loans")
)

type LoanUsecase struct {
	loanRepo       repository.LoanRepository
	userUsecase    UserUsecaseInterface
	paymentUsecase *PaymentUsecase
}

func NewLoanUsecase(
	loanRepo repository.LoanRepository,
	userUsecase UserUsecaseInterface,
	paymentUsecase *PaymentUsecase,
) *LoanUsecase {
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
	if _, err := u.userUsecase.IsUserDelinquent(ctx, loan.UserID); err != nil {
		return err
	}

	// only
	// activeStatus := entity.LoanStatusActive
	// if loans, _ := lu.GetLoansByUserID(ctx, loan.UserID, activeStatus); loans != nil {
	// 	return ErrStillHasActiveLoan
	// }

	return nil

}

func (u *LoanUsecase) CreateLoanWithPayments(ctx context.Context, loan *entity.Loan) error {

	if err := u.validateBillingStartDate(loan.BillingStartDate); err != nil {
		return err
	}

	if _, err := u.userUsecase.GetByID(ctx, loan.UserID); err != nil {
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

	loan, err = u.loanRepo.CreateLoanInTx(tx, loan)
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

	err = u.paymentUsecase.CreatePaymentsInTx(tx, paymentsPayload)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (lu *LoanUsecase) validateBillingStartDate(billingStartDate time.Time) error {
	// if billingStartDate.Before(time.Now().Truncate(24 * time.Hour)) {
	// 	return ErrInvalidBillingStartDate
	// }
	return nil
}

func (lu *LoanUsecase) calculateTotalOutstanding(loan *entity.Loan) float64 {
	totalInterest := lu.calculateInterest(loan)

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
