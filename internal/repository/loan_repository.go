package repository

import (
	"context"
	"database/sql"
	"errors"
	"loan-management/internal/entity"
	"time"
)

var (
	ErrLoanNotFound = errors.New("loan not found")
)

type LoanRepository interface {
	Create(ctx context.Context, loan *entity.Loan) error
	GetByID(ctx context.Context, id int64) (*entity.Loan, error)
	GetAll(ctx context.Context) ([]*entity.Loan, error)
	GetByUserID(ctx context.Context, userId int64, status *entity.LoanStatus) ([]*entity.Loan, error)
}

type loanRepository struct {
	db *sql.DB
}

func NewLoanRepository(db *sql.DB) LoanRepository {
	return &loanRepository{db: db}
}

func (r *loanRepository) Create(ctx context.Context, loan *entity.Loan) error {
	query := `
		INSERT INTO loans (
			user_id,
			interest,
			interest_type,
			tenure,
			tenure_type,
			amount,
			outstanding,
			status,
			created_at,
			billing_start_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	result, err := r.db.ExecContext(ctx, query,
		loan.UserID,
		loan.Interest,
		loan.InterestType,
		loan.Tenure,
		loan.TenureType,
		loan.Amount,
		loan.Outstanding,
		loan.Status,
		time.Now(),
		loan.BillingStartDate,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	loan.ID = id
	return nil
}

func scanLoan(scanner interface{ Scan(dest ...any) error }, loan *entity.Loan) error {
	return scanner.Scan(
		&loan.ID,
		&loan.UserID,
		&loan.Interest,
		&loan.InterestType,
		&loan.Tenure,
		&loan.TenureType,
		&loan.Amount,
		&loan.Outstanding,
		&loan.Status,
		&loan.CreatedAt,
		&loan.BillingStartDate,
	)
}

func (r *loanRepository) GetAll(ctx context.Context) ([]*entity.Loan, error) {
	query := `SELECT * FROM loans`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loans []*entity.Loan
	for rows.Next() {
		loan := entity.Loan{}
		if err := scanLoan(rows, &loan); err != nil {
			return nil, err
		}
		loans = append(loans, &loan)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return loans, nil
}

func (r *loanRepository) GetByID(ctx context.Context, id int64) (*entity.Loan, error) {
	query := `SELECT * FROM loans WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query)
	loan := entity.Loan{}
	if err := scanLoan(row, &loan); err != nil {
		return nil, err
	}

	return &loan, nil
}

func (r *loanRepository) GetByUserID(ctx context.Context, userID int64, status *entity.LoanStatus) ([]*entity.Loan, error) {
	query := `SELECT * FROM loans WHERE user_id = ?`
	args := []interface{}{userID}

	if status != nil {
		query += ` AND status = ?`
		args = append(args, *status)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loans []*entity.Loan
	for rows.Next() {
		loan := entity.Loan{}
		if err := scanLoan(rows, &loan); err != nil {
			return nil, err
		}
		loans = append(loans, &loan)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return loans, nil
}
