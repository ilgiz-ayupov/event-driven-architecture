package auth_usecase

import (
	"event-driven-architecture/internal/domain"
	"event-driven-architecture/internal/usecase"
)

type AuthUserUsecase struct {
	userRepo       usecase.UserRepo
	passwordHasher usecase.PasswordHasher
}

func NewAuthUser(
	userRepo usecase.UserRepo,
	passwordHasher usecase.PasswordHasher,
) *AuthUserUsecase {
	return &AuthUserUsecase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
	}
}

func (u *AuthUserUsecase) Execute(ctx usecase.AppCtx, in AuthUserInput) (domain.AuthUser, error) {
	var zero domain.AuthUser

	// поиск пользователя
	user, err := u.userRepo.FindByEmail(ctx, in.Email)
	switch err {
	case nil:
	case usecase.ErrNoData:
		return zero, usecase.ErrInvalidCredentials
	default:
		return zero, err
	}

	// проверка пароля
	switch ok, err := u.passwordHasher.Compare(in.Password, user.PasswordHash); err {
	case nil:
		if !ok {
			return zero, usecase.ErrInvalidCredentials
		}
	default:
		return zero, err
	}

	return domain.NewAuthUser(user.ID), nil
}
