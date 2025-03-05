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
	CreatePayment(ctx context.Context, loan *entity.Payment) error
	GetPaymentByID(ctx context.Context, id int64) (*entity.Payment, error)
	GetAllPayments(ctx context.Context, status *entity.PaymentStatus) ([]*entity.Payment, error)
	GetPaymentsByLoanID(ctx context.Context, loanId int64, status *entity.PaymentStatus) ([]*entity.Payment, error)
}

type paymentRepository struct {
	db *sql.DB
}

func (r *paymentRepository) NewPaymentRepository(db *sql.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) CreatePayment(ctx context.Context, payment *entity.Payment) error {
	query := `
		INSERT INTO payments (
			loan_id,
			due_date,
			amount,
			interest,
			total_amount,
			status,
			paid_at,
			created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?);
	`

	result, err := r.db.ExecContext(ctx, query,
		payment.LoanID,
		payment.DueDate,
		payment.Amount,
		payment.Status,
		payment.PaidAt,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	payment.ID = id
	return nil
}

func scanPayment(scanner interface{ Scan(dest ...any) error }, payment *entity.Payment) error {
	var paidAt sql.NullTime

	res := scanner.Scan(
		&payment.ID,
		&payment.LoanID,
		&payment.DueDate,
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
		SELECT id, loan_id, due_date, amount, interest, total_amount, status, paid_at, created_at
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
		SELECT id, loan_id, due_date, amount, interest, total_amount, status, paid_at, created_at
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

func (r *paymentRepository) GetPaymentsByLoanID(ctx context.Context, loanId int64, status *entity.PaymentStatus) ([]*entity.Payment, error) {
	query := `
		SELECT id, loan_id, due_date, amount, interest, total_amount, status, paid_at, created_at
		FROM payments
		WHERE loan_id = ?
	`
	args := []interface{}{loanId}

	if status != nil {
		query += ` AND status = ?`
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
