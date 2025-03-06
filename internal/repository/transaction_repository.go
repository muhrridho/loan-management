package repository

import (
	"context"
	"database/sql"
	"loan-management/internal/entity"
)

type transactionRepository struct {
	db *sql.DB
}

type TransactionRepository interface {
	CreateTransaction(tx *sql.Tx, transaction *entity.Transaction) (int64, error)
	GetTransactionByID(ctx context.Context, id int64) (*entity.Transaction, error)
	BeginTx() (*sql.Tx, error)
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateTransaction(tx *sql.Tx, transaction *entity.Transaction) (int64, error) {
	query := `
	INSERT INTO transactions (
		total_amount,
		penalty,
		status,
		paid_at,
		created_at
	) VALUES (?, ?, ?, ?, ?)
	`

	result, err := tx.Exec(
		query,
		transaction.TotalAmount,
		transaction.Penalty,
		transaction.Status,
		transaction.PaidAt,
		transaction.CreatedAt,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (r *transactionRepository) GetTransactionByID(ctx context.Context, id int64) (*entity.Transaction, error) {
	query := `
	SELECT id, total_amount, penalty, status, paid_at, created_at
	FROM transactions
	WHERE id = ?
	`
	row := r.db.QueryRowContext(ctx, query, id)

	transaction := &entity.Transaction{}
	var paidAt sql.NullTime

	err := row.Scan(&transaction.ID, &transaction.TotalAmount, &transaction.Penalty, &transaction.Status, &paidAt, &transaction.CreatedAt)
	if err != nil {
		return nil, err
	}

	if paidAt.Valid {
		transaction.PaidAt = &paidAt.Time
	}

	return transaction, nil
}

func (r *transactionRepository) BeginTx() (*sql.Tx, error) {
	return r.db.Begin()
}
