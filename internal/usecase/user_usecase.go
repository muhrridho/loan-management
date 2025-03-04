package usecase

import (
	"context"
	"errors"
	"loan-management/internal/entity"
	"loan-management/internal/repository"
)

var (
	ErrEmailAlreadyUsed     = errors.New("Your email is already being used")
	ErrMissingRequiredField = errors.New("Name & Email is required")
)

type UserUsecaseInterface interface {
	RegisterUser(ctx context.Context, user *entity.User) error
	GetAll(ctx context.Context) ([]*entity.User, error)
	GetByID(ctx context.Context, id int64) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	IsUserDelinquent(ctx context.Context, userID int64) (bool, error)
}

type UserUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) *UserUsecase {
	return &UserUsecase{userRepo: userRepo}
}

func (uc *UserUsecase) RegisterUser(ctx context.Context, user *entity.User) error {
	if user.Email == "" || user.Name == "" {
		return ErrMissingRequiredField
	}

	if u, _ := uc.GetByEmail(ctx, user.Email); u != nil {
		return ErrEmailAlreadyUsed
	}

	return uc.userRepo.Create(ctx, user)
}

func (uc *UserUsecase) GetAll(ctx context.Context) ([]*entity.User, error) {
	return uc.userRepo.GetAll(ctx)
}

func (uc *UserUsecase) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	return uc.userRepo.GetByID(ctx, id)
}

func (uc *UserUsecase) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	return uc.userRepo.GetByEmail(ctx, email)
}

func (uc *UserUsecase) IsUserDelinquent(ctx context.Context, userID int64) (bool, error) {
	if ctx == nil || userID == 0 {
		return false, errors.New("awsdfasd")
	}
	return false, nil
}
