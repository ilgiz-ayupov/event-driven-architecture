package user

import (
	"context"
	"event-driven-architecture/internal/usecase"
)

type UserService struct {
	tsManager     usecase.TransactionManager
	appCtxManager usecase.AppCtxManager

	createUserUseCase *usecase.CreateUserUseCase
}

func NewUserService(
	tsManager usecase.TransactionManager,
	appCtxManager usecase.AppCtxManager,
	createUserUseCase *usecase.CreateUserUseCase,
) *UserService {
	return &UserService{
		tsManager:         tsManager,
		appCtxManager:     appCtxManager,
		createUserUseCase: createUserUseCase,
	}
}

func (s *UserService) CreateUser(ctx context.Context, name string, email string) (int64, error) {
}
