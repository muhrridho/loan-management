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
	CreateLoan(ctx context.Context, loan *entity.Loan) error
	CreateLoanInTx(tx *sql.Tx, loan *entity.Loan) (*entity.Loan, error)
	GetLoanByID(ctx context.Context, id int64, status *entity.LoanStatus) (*entity.Loan, error)
	GetAllLoans(ctx context.Context) ([]*entity.Loan, error)
	GetLoansByUserID(ctx context.Context, userId int64, status *entity.LoanStatus) ([]*entity.Loan, error)
	UpdateLoanOutstanding(tx *sql.Tx, outstanding float64, loanID int64) error
	BeginTx() (*sql.Tx, error)
}

// type loanRepository struct {
// 	db *sql.DB
// }

//	func NewLoanRepository(db *sql.DB) LoanRepository {
//		return &loanRepository{db: db}
//	}
type loanRepository struct {
	db *sql.DB
}

func NewLoanRepository(db *sql.DB) LoanRepository {
	return &loanRepository{db: db}
}

func (r *loanRepository) CreateLoan(ctx context.Context, loan *entity.Loan) error {
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

func (r *loanRepository) CreateLoanInTx(tx *sql.Tx, loan *entity.Loan) (*entity.Loan, error) {
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
	result, err := tx.Exec(query,
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
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	loan.ID = id

	return loan, nil
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

func (r *loanRepository) GetAllLoans(ctx context.Context) ([]*entity.Loan, error) {
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

func (r *loanRepository) GetLoanByID(ctx context.Context, id int64, status *entity.LoanStatus) (*entity.Loan, error) {
	query := `SELECT * FROM loans WHERE id = ?`
	args := []interface{}{id}

	if status != nil {
		query += ` AND status = ?`
		args = append(args, *status)
	}

	row := r.db.QueryRowContext(ctx, query, args...)
	loan := entity.Loan{}
	if err := scanLoan(row, &loan); err != nil {
		return nil, err
	}

	return &loan, nil
}

func (r *loanRepository) GetLoansByUserID(ctx context.Context, userID int64, status *entity.LoanStatus) ([]*entity.Loan, error) {
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

func (r *loanRepository) UpdateLoanOutstanding(tx *sql.Tx, outstanding float64, loanID int64) error {
	var query string

	if outstanding == 0 {
		query = `UPDATE loans SET outstanding = ?, status = 99 WHERE id = ?`
	} else {
		query = `UPDATE loans SET outstanding = ? WHERE id = ?`
	}

	_, err := tx.Exec(query, outstanding, loanID)
	if err != nil {
		return err
	}

	return nil

}

func (r *loanRepository) BeginTx() (*sql.Tx, error) {
	return r.db.Begin()
}
