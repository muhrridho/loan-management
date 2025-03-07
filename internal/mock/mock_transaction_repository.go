package mock

import (
	"context"
	"database/sql"
	"loan-management/internal/entity"

	"github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) CreateTransaction(tx *sql.Tx, transaction *entity.Transaction) (int64, error) {
	args := m.Called(tx, transaction)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockTransactionRepository) GetTransactionByID(ctx context.Context, id int64) (*entity.Transaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.Transaction), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTransactionRepository) BeginTx() (*sql.Tx, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).(*sql.Tx), args.Error(1)
	}
	return nil, args.Error(1)
}
