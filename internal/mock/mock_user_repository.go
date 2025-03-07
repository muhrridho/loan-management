package mock

import (
	"context"
	"loan-management/internal/entity"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetAllUsers(ctx context.Context) ([]*entity.User, error) {
	args := m.Called(ctx)
	if users, ok := args.Get(0).([]*entity.User); ok {
		return users, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	args := m.Called(ctx, id)
	if user, ok := args.Get(0).(*entity.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if user, ok := args.Get(0).(*entity.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}
