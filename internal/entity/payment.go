package entity

import "time"

type PaymentStatus int8

const (
	PaymentStatusActive PaymentStatus = iota
	PaymentStatusPaid
)

type Payment struct {
	ID          int64         `db:"id"`
	LoanID      int64         `db:"loan_id"`
	DueDate     time.Time     `db:"due_date"`
	Amount      float64       `db:"amount"`
	Interest    float64       `db:"interest"`
	TotalAmount float64       `db:"total_amount"`
	Status      PaymentStatus `db:"status"`
	PaidAt      *time.Time    `db:"paid_at"`
	CreatedAt   time.Time     `db:"created_at"`
}

type CreatePaymentPayload struct {
	LoanID      int64     `json:"loan_id"`
	DueDate     time.Time `json:"due_date"`
	Amount      float64   `json:"amount"`
	Interest    float64   `json:"interest"`
	TotalAmount float64   `json:"total_amount"`
}
