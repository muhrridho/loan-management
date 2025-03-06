package usecase

import (
	"context"
	"errors"
	"loan-management/internal/entity"
	"loan-management/internal/repository"
)

var (
	ErrUserNotFound         = errors.New("User not found")
	ErrEmailAlreadyUsed     = errors.New("Your email is already being used")
	ErrMissingRequiredField = errors.New("Name & Email is required")
)

type UserUsecaseInterface interface {
	RegisterUser(ctx context.Context, user *entity.User) error
	GetAllUsers(ctx context.Context) ([]*entity.User, error)
	GetUserByID(ctx context.Context, id int64) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	IsUserDelinquent(ctx context.Context, userID int64) (bool, error)
}

type UserUsecase struct {
	userRepo    repository.UserRepository
	loanUsecase LoanUsecase
}

func NewUserUsecase(userRepo repository.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

func (u *UserUsecase) InjectDependencies(loanUsecase *LoanUsecase) {
	u.loanUsecase = *loanUsecase
}

func (u *UserUsecase) RegisterUser(ctx context.Context, user *entity.User) error {
	if user.Email == "" || user.Name == "" {
		return ErrMissingRequiredField
	}

	if user, _ := u.GetUserByEmail(ctx, user.Email); user != nil {
		return ErrEmailAlreadyUsed
	}

	return u.userRepo.CreateUser(ctx, user)
}

func (u *UserUsecase) GetAllUsers(ctx context.Context) ([]*entity.User, error) {
	return u.userRepo.GetAllUsers(ctx)
}

func (u *UserUsecase) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	return u.userRepo.GetUserByID(ctx, id)
}

func (u *UserUsecase) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return u.userRepo.GetUserByEmail(ctx, email)
}

func (u *UserUsecase) IsUserDelinquent(ctx context.Context, userID int64) (bool, error) {

	loanStatusActive := entity.LoanStatusActive
	activeLoans, err := u.loanUsecase.GetLoansByUserID(ctx, userID, loanStatusActive)

	if err != nil {
		return false, err
	}

	if len(activeLoans) == 0 {
		return false, nil
	}

	for _, loan := range activeLoans {
		if numOfDuePayments, _ := u.loanUsecase.GetLoanDuePayments(ctx, loan); len(numOfDuePayments) > 2 {
			return true, nil
		}
	}

	return false, nil
}
