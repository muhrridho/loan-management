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
)

type LoanUsecase struct {
	loanRepo repository.LoanRepository
}

func NewLoanUsecase(loanRepo repository.LoanRepository) *LoanUsecase {
	return &LoanUsecase{loanRepo: loanRepo}
}

func (lu *LoanUsecase) Create(ctx context.Context, loan *entity.Loan) error {
	fmt.Println(loan)
	if err := lu.validateBillingStartDate(loan.BillingStartDate); err != nil {
		return err
	}

	if loan.InterestType == entity.InterestTypeFlatAnnual {

		loan.Outstanding = lu.calculateTotalOutstanding(loan)
	} else {
		// TODO: Implement Reduce Annual Type
	}
	return lu.loanRepo.Create(ctx, loan)
}

func (lu *LoanUsecase) GetAll(ctx context.Context) ([]*entity.Loan, error) {
	return lu.loanRepo.GetAll(ctx)
}

func (lu *LoanUsecase) GetByID(ctx context.Context, id int64) (*entity.Loan, error) {
	return lu.loanRepo.GetByID(ctx, id)
}

func (lu *LoanUsecase) validateBillingStartDate(billingStartDate time.Time) error {
	fmt.Println(billingStartDate)
	if billingStartDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return ErrInvalidBillingStartDate
	}
	return nil
}

func (lu LoanUsecase) calculateTotalOutstanding(loan *entity.Loan) float64 {
	tenureInYears := float64(loan.Tenure) / 52
	if loan.TenureType == entity.TenureTypeMonthly {
		tenureInYears = float64(loan.Tenure) / 12
	}

	totalInterest := loan.Amount * (loan.Interest / 100) * tenureInYears

	totalOutstanding := loan.Amount + totalInterest

	return totalOutstanding
}
