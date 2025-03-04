package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"loan-management/internal/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetAll(ctx context.Context) ([]*entity.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

var mockUser = &entity.User{
	Email:     "test@test",
	Name:      "test",
	CreatedAt: time.Now(),
}

func TestRegisterUser(t *testing.T) {
	t.Run("Success User Register", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		userUsecase := NewUserUsecase(mockRepo)

		mockRepo.On("GetByEmail", mock.Anything, mockUser.Email).Return(nil, nil)
		mockRepo.On("Create", mock.Anything, mockUser).Return(nil)

		err := userUsecase.RegisterUser(context.Background(), mockUser)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed User Refister", func(t *testing.T) {
		mockRepo := new(MockUserRepository)

		userUsecase := NewUserUsecase(mockRepo)

		expectedError := errors.New("error")
		mockRepo.On("GetByEmail", mock.Anything, mockUser.Email).Return(nil, nil)
		mockRepo.On("Create", mock.Anything, mockUser).Return(expectedError)

		err := userUsecase.RegisterUser(context.Background(), mockUser)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed User Register - Missing field", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		userUsecase := NewUserUsecase(mockRepo)

		customMockUser := *mockUser
		customMockUser.Name = ""

		mockRepo.On("GetByEmail", mock.Anything, customMockUser.Email).Return(nil, nil)
		err := userUsecase.RegisterUser(context.Background(), &customMockUser)

		assert.Error(t, err)
		assert.Equal(t, "Name & Email is required", err.Error())

		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Failed User Register - Email already used", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		userUsecase := NewUserUsecase(mockRepo)

		mockRepo.On("GetByEmail", mock.Anything, mockUser.Email).Return(mockUser, nil)

		err := userUsecase.RegisterUser(context.Background(), mockUser)

		assert.Error(t, err)
		assert.Equal(t, "Your email is already being used", err.Error())

		mockRepo.AssertNotCalled(t, "Create")
	})
}

func TestUserGetAll(t *testing.T) {
	t.Run("Success Get All Users", func(t *testing.T) {
		mockRepo := new(MockUserRepository)

		userUsecase := NewUserUsecase(mockRepo)

		expectedUsers := []*entity.User{mockUser}

		mockRepo.On("GetAll", mock.Anything).Return(expectedUsers, nil)
		users, err := userUsecase.GetAll(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		assert.Len(t, users, 1)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Get All Users", func(t *testing.T) {
		mockRepo := new(MockUserRepository)

		userUsecase := NewUserUsecase(mockRepo)

		expectedError := errors.New("database error")
		mockRepo.On("GetAll", mock.Anything).Return([]*entity.User{}, expectedError)

		users, err := userUsecase.GetAll(context.Background())

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Empty(t, users)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserGetByID(t *testing.T) {
	t.Run("Success Get User by ID", func(t *testing.T) {
		mockRepo := new(MockUserRepository)

		userUsecase := NewUserUsecase(mockRepo)

		expectedUser := mockUser
		expectedUser.ID = 1

		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(expectedUser, nil)

		user, err := userUsecase.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Get User by ID", func(t *testing.T) {
		mockRepo := new(MockUserRepository)

		userUsecase := NewUserUsecase(mockRepo)

		expectedError := errors.New("user not found")
		mockRepo.On("GetByID", mock.Anything, int64(69)).Return((*entity.User)(nil), expectedError)

		user, err := userUsecase.GetByID(context.Background(), 69)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}
