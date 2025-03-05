package entity

import (
	"fmt"
	"time"
)

type LoanStatus int8

const (
	LoanStatusActive LoanStatus = 1
	LoanStatusPaid   LoanStatus = 99
)

func (it LoanStatus) String() string {
	switch it {
	case LoanStatusActive:
		return "Active"
	case LoanStatusPaid:
		return "Paid"
	default:
		return "Unknown"
	}
}

type InterestType int8

const (
	InterestTypeFlatAnnual InterestType = iota
	InterestTypeReducingAnnual
)

func (it InterestType) String() string {
	switch it {
	case InterestTypeFlatAnnual:
		return "Flat Annual"
	case InterestTypeReducingAnnual:
		return "Reducing Annual"
	default:
		return "Unknown"
	}
}

type TenureType int8

const (
	TenureTypeWeekly TenureType = iota
	// TenureTypeMonthly
)

func (it TenureType) String() string {
	switch it {
	case TenureTypeWeekly:
		return "Weeks"
	// case TenureTypeMonthly:
	// 	return "Months"
	default:
		return "Unknown"
	}
}

type Loan struct {
	ID               int64        `db:"id"`
	UserID           int64        `db:"user_id"`
	Interest         float64      `db:"interest"`
	InterestType     InterestType `db:"interest_type"`
	Tenure           int          `db:"tenure"`
	TenureType       TenureType   `db:"tenure_type"`
	Amount           float64      `db:"amount"`
	Outstanding      float64      `db:"outstanding"`
	Status           LoanStatus   `db:"status"`
	CreatedAt        time.Time    `db:"created_at"`
	BillingStartDate time.Time    `db:"billing_start_date"`
}

func (l Loan) String() string {
	return fmt.Sprintf(
		"Loan ID: %d\n"+
			"User ID: %d\n"+
			"Amount: %.2f\n"+
			"Outstanding: %.2f\n"+
			"Interest: %.2f%%\n"+
			"Interest Type: %s\n"+
			"Tenure: %d %s\n"+
			"Status: %s\n"+
			"Created At: %s\n"+
			"Billing Start Date: %s\n",
		l.ID,
		l.UserID,
		l.Amount,
		l.Outstanding,
		l.Interest,
		l.InterestType,
		l.Tenure,
		l.TenureType,
		l.Status,
		l.CreatedAt.Format("2006-01-02 15:04:05"),
		l.BillingStartDate.Format("2006-01-02"),
	)
}

type CreateLoanPayload struct {
	UserID           int64        `json:"user_id"`
	Amount           float64      `json:"amount"`
	Interest         float64      `json:"interest"`
	InterestType     InterestType `json:"interest_type"`
	Tenure           int          `json:"tenure"`
	TenureType       TenureType   `json:"tenure_type"`
	BillingStartDate time.Time    `json:"billing_start_date"`
}

func NewLoan(userID int64, amount float64, interest float64, tenure int, interestType InterestType, tenureType TenureType, billingStartDate time.Time) *Loan {
	outstanding := amount + (amount * float64(interest))

	return &Loan{
		UserID:           userID,
		Interest:         interest,
		InterestType:     interestType,
		Tenure:           tenure,
		TenureType:       tenureType,
		Amount:           amount,
		Outstanding:      outstanding,
		Status:           LoanStatusActive,
		CreatedAt:        time.Now(),
		BillingStartDate: billingStartDate,
	}
}
