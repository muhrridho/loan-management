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
	GetAll(ctx context.Context) ([]*entity.User, error)
	GetByID(ctx context.Context, id int64) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
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

	loanStatusActive := entity.LoanStatusActive
	activeLoans, err := uc.loanUsecase.GetLoansByUserID(ctx, userID, loanStatusActive)

	if err != nil {
		return false, err
	}

	if len(activeLoans) == 0 {
		return false, nil
	}

	for _, loan := range activeLoans {
		if numOfDuePayments, _ := uc.loanUsecase.GetLoanDuePayments(ctx, loan); len(numOfDuePayments) > 2 {
			return true, nil
		}
	}

	return false, nil
}
