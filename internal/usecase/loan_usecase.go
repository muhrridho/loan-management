package usecase

import (
	"context"
	"errors"
	"fmt"
	"loan-management/internal/entity"
	"loan-management/internal/repository"
	"time"
)

var (
	ErrInvalidBillingStartDate = errors.New("billing start date cannot be in the past")
	ErrStillHasActiveLoan      = errors.New("Can't create loan because you still have an active loans")
)

type LoanUsecase struct {
	loanRepo    repository.LoanRepository
	userUsecase UserUsecaseInterface
}

func NewLoanUsecase(loanRepo repository.LoanRepository, userUsecase UserUsecaseInterface) *LoanUsecase {
	return &LoanUsecase{
		loanRepo:    loanRepo,
		userUsecase: userUsecase,
	}
}

func (lu *LoanUsecase) CreateLoan(ctx context.Context, loan *entity.Loan) error {
	if err := lu.validateBillingStartDate(loan.BillingStartDate); err != nil {
		return err
	}

	if _, err := lu.userUsecase.GetByID(ctx, loan.UserID); err != nil {
		return err
	}

	if err := lu.CheckCreateLoanEligibility(ctx, loan); err != nil {
		return err
	}

	if loan.InterestType == entity.InterestTypeFlatAnnual {
		loan.Outstanding = lu.calculateTotalOutstanding(loan)
	} else {
		// TODO: Implement Reduce Annual Type
	}
	return lu.loanRepo.CreateLoan(ctx, loan)
}

func (lu *LoanUsecase) GetAllLoans(ctx context.Context) ([]*entity.Loan, error) {
	return lu.loanRepo.GetAllLoans(ctx)
}

func (lu *LoanUsecase) GetLoanByID(ctx context.Context, id int64) (*entity.Loan, error) {
	return lu.loanRepo.GetLoanByID(ctx, id)
}

func (lu *LoanUsecase) GetLoansByUserID(ctx context.Context, userID int64, status entity.LoanStatus) ([]*entity.Loan, error) {
	return lu.loanRepo.GetLoansByUserID(ctx, userID, &status)
}

func (lu *LoanUsecase) CheckCreateLoanEligibility(ctx context.Context, loan *entity.Loan) error {
	if _, err := lu.userUsecase.IsUserDelinquent(ctx, loan.UserID); err != nil {
		return err
	}

	activeStatus := entity.LoanStatusActive
	if loans, _ := lu.GetLoansByUserID(ctx, loan.UserID, activeStatus); loans != nil {
		return ErrStillHasActiveLoan
	}

	return nil

}

func (lu *LoanUsecase) validateBillingStartDate(billingStartDate time.Time) error {
	fmt.Println(billingStartDate)
	if billingStartDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return ErrInvalidBillingStartDate
	}
	return nil
}

func (lu *LoanUsecase) calculateTotalOutstanding(loan *entity.Loan) float64 {
	tenureInYears := float64(loan.Tenure) / 52
	if loan.TenureType == entity.TenureTypeMonthly {
		tenureInYears = float64(loan.Tenure) / 12
	}

	totalInterest := loan.Amount * (loan.Interest / 100) * tenureInYears

	totalOutstanding := loan.Amount + totalInterest

	return totalOutstanding
}
