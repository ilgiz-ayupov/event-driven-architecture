package user_usecase

import (
	"event-driven-architecture/internal/domain"
	"event-driven-architecture/internal/usecase"
)

type CreateUserUseCase struct {
	idGenerator    usecase.IDGenerator
	passwordHasher usecase.PasswordHasher

	userRepo usecase.UserRepo
}

func NewCreateUser(
	idGenerator usecase.IDGenerator,
	passwordHasher usecase.PasswordHasher,

	userRepo usecase.UserRepo,
) *CreateUserUseCase {
	return &CreateUserUseCase{
		idGenerator:    idGenerator,
		passwordHasher: passwordHasher,

		userRepo: userRepo,
	}
}

func (u *CreateUserUseCase) Execute(ctx usecase.AppCtx, in CreateUserInput) (CreateUserOutput, error) {
	var zero CreateUserOutput

	// хешировать пароль
	hash, err := u.passwordHasher.Hash(in.Password)
	if err != nil {
		return zero, err
	}

	// сохранить пользователя
	user := domain.NewUser(
		u.idGenerator.NewID(),
		in.Email,
		hash,
	)

	if err := u.userRepo.Create(ctx, user); err != nil {
		return zero, err
	}

	return NewCreateUserOutput(user.ID), nil
}
