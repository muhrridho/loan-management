package repository

import (
	"context"
	"database/sql"
	"errors"
	"loan-management/internal/entity"
	"time"
)

var (
	ErrPaymentNotFound = errors.New("loan not found")
)

type PaymentRepository interface {
	CreatePaymentsWithTx(tx *sql.Tx, payments []*entity.Payment) error
	GetPaymentByID(ctx context.Context, id int64) (*entity.Payment, error)
	GetAllPayments(ctx context.Context, status *entity.PaymentStatus) ([]*entity.Payment, error)
	GetPaymentsByLoanID(ctx context.Context, loanId int64, status *entity.PaymentStatus, dueBefore *time.Time) ([]*entity.Payment, error)
}

type paymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) CreatePaymentsWithTx(tx *sql.Tx, payments []*entity.Payment) error {
	if len(payments) == 0 {
		return errors.New("no payments to create")
	}

	query := `
		INSERT INTO payments (
			loan_id,
			due_date,
			payment_no,
			amount,
			interest,
			total_amount,
			status,
			paid_at,
			created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, payment := range payments {
		_, err = stmt.Exec(
			payment.LoanID,
			payment.DueDate,
			payment.PaymentNo,
			payment.Amount,
			payment.Interest,
			payment.TotalAmount,
			payment.Status,
			payment.PaidAt,
			payment.CreatedAt,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func scanPayment(scanner interface{ Scan(dest ...any) error }, payment *entity.Payment) error {
	var paidAt sql.NullTime

	res := scanner.Scan(
		&payment.ID,
		&payment.LoanID,
		&payment.DueDate,
		&payment.PaymentNo,
		&payment.Amount,
		&payment.Interest,
		&payment.TotalAmount,
		&payment.Status,
		&paidAt,
		&payment.CreatedAt,
	)

	if paidAt.Valid {
		payment.PaidAt = &paidAt.Time
	}

	return res

}

func (r *paymentRepository) GetPaymentByID(ctx context.Context, id int64) (*entity.Payment, error) {
	query := `
		SELECT id, loan_id, due_date, payment_no, amount, interest, total_amount, status, paid_at, created_at
		FROM payments
		WHERE id = ?
	`

	payment := &entity.Payment{}

	err := scanPayment(r.db.QueryRowContext(ctx, query, id), payment)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPaymentNotFound
		}
		return nil, err
	}

	return payment, nil
}

func (r *paymentRepository) GetAllPayments(ctx context.Context, status *entity.PaymentStatus) ([]*entity.Payment, error) {
	query := `
	SELECT id, loan_id, due_date, payment_no, amount, interest, total_amount, status, paid_at, created_at
	FROM payments
	`
	args := []interface{}{}

	if status != nil {
		query += ` WHERE status = ?`
		args = append(args, *status)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*entity.Payment
	for rows.Next() {
		payment := entity.Payment{}
		if err := scanPayment(rows, &payment); err != nil {
			return nil, err
		}
		payments = append(payments, &payment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}

func (r *paymentRepository) GetPaymentsByLoanID(ctx context.Context, loanId int64, status *entity.PaymentStatus, dueBefore *time.Time) ([]*entity.Payment, error) {
	query := `
		SELECT id, loan_id, due_date, payment_no, amount, interest, total_amount, status, paid_at, created_at
		FROM payments
		WHERE loan_id = ?
	`
	args := []interface{}{loanId}

	if status != nil {
		query += ` AND status = ?`
		args = append(args, *status)
	}

	if dueBefore != nil {
		query += ` AND due_date <= ?`
		args = append(args, dueBefore)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*entity.Payment
	for rows.Next() {
		payment := entity.Payment{}
		if err := scanPayment(rows, &payment); err != nil {
			return nil, err
		}
		payments = append(payments, &payment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}
