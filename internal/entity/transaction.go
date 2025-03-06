package entity

import "time"

type TransactionStatus int8

const (
	TransactionStatusActive TransactionStatus = 1
	TransactionStatusPaid   TransactionStatus = 99
)

type TransactionInquiry struct {
	LoanID     int64      `json:"loan_id"`
	AmountDue  float64    `json:"amount_due"`
	DueDate    time.Time  `json:"due_date"`
	LoanDetail *Loan      `json:"loan_detail"`
	Bills      []*Payment `json:"payments"`
}

type Transaction struct {
	ID          int64             `db:"id"`
	TotalAmount float64           `db:"total_amount"`
	Penalty     float64           `db:"penalty"`
	Status      TransactionStatus `db:"status"`
	PaidAt      *time.Time        `db:"paid_at"`
	CreatedAt   time.Time         `db:"created_at"`
}

type CreateTransactionPayload struct {
	LoanID int64   `json:"loan_id"`
	Amount float64 `json:"amount"`
}
