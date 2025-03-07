package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"loan-management/internal/entity"

	internalMock "loan-management/internal/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var MockUser = &entity.User{
	Email:     "test@test",
	Name:      "test",
	CreatedAt: time.Now(),
}

func TestRegisterUser(t *testing.T) {
	t.Run("Success User Register", func(t *testing.T) {
		mockRepo := new(internalMock.MockUserRepository)
		userUsecase := NewUserUsecase(mockRepo)

		mockRepo.On("GetUserByEmail", mock.Anything, MockUser.Email).Return(nil, nil)
		mockRepo.On("CreateUser", mock.Anything, MockUser).Return(nil)

		err := userUsecase.RegisterUser(context.Background(), MockUser)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed User Refister", func(t *testing.T) {
		mockRepo := new(internalMock.MockUserRepository)

		userUsecase := NewUserUsecase(mockRepo)

		expectedError := errors.New("error")
		mockRepo.On("GetUserByEmail", mock.Anything, MockUser.Email).Return(nil, nil)
		mockRepo.On("CreateUser", mock.Anything, MockUser).Return(expectedError)

		err := userUsecase.RegisterUser(context.Background(), MockUser)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed User Register - Missing field", func(t *testing.T) {
		mockRepo := new(internalMock.MockUserRepository)
		userUsecase := NewUserUsecase(mockRepo)

		customMockUser := *MockUser
		customMockUser.Name = ""

		mockRepo.On("GetUserByEmail", mock.Anything, customMockUser.Email).Return(nil, nil)
		err := userUsecase.RegisterUser(context.Background(), &customMockUser)

		assert.Error(t, err)
		assert.Equal(t, "Name & Email is required", err.Error())

		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Failed User Register - Email already used", func(t *testing.T) {
		mockRepo := new(internalMock.MockUserRepository)
		userUsecase := NewUserUsecase(mockRepo)

		mockRepo.On("GetUserByEmail", mock.Anything, MockUser.Email).Return(MockUser, nil)

		err := userUsecase.RegisterUser(context.Background(), MockUser)

		assert.Error(t, err)
		assert.Equal(t, "Your email is already being used", err.Error())

		mockRepo.AssertNotCalled(t, "Create")
	})
}

func TestUserGetAllUsers(t *testing.T) {
	t.Run("Success Get All Users", func(t *testing.T) {
		mockRepo := new(internalMock.MockUserRepository)

		userUsecase := NewUserUsecase(mockRepo)

		expectedUsers := []*entity.User{MockUser}

		mockRepo.On("GetAllUsers", mock.Anything).Return(expectedUsers, nil)
		users, err := userUsecase.GetAllUsers(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		assert.Len(t, users, 1)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Get All Users", func(t *testing.T) {
		mockRepo := new(internalMock.MockUserRepository)

		userUsecase := NewUserUsecase(mockRepo)

		expectedError := errors.New("database error")
		mockRepo.On("GetAllUsers", mock.Anything).Return([]*entity.User{}, expectedError)

		users, err := userUsecase.GetAllUsers(context.Background())

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Empty(t, users)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserGetByID(t *testing.T) {
	t.Run("Success Get User by ID", func(t *testing.T) {
		mockRepo := new(internalMock.MockUserRepository)

		userUsecase := NewUserUsecase(mockRepo)

		expectedUser := MockUser
		expectedUser.ID = 1

		mockRepo.On("GetUserByID", mock.Anything, int64(1)).Return(expectedUser, nil)

		user, err := userUsecase.GetUserByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Get User by ID", func(t *testing.T) {
		mockRepo := new(internalMock.MockUserRepository)

		userUsecase := NewUserUsecase(mockRepo)

		expectedError := errors.New("user not found")
		mockRepo.On("GetUserByID", mock.Anything, int64(69)).Return((*entity.User)(nil), expectedError)

		user, err := userUsecase.GetUserByID(context.Background(), 69)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}
