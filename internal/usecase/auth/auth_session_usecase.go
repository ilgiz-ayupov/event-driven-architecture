package auth_usecase

import (
	"event-driven-architecture/internal/domain"
	"event-driven-architecture/internal/usecase"
)

type AuthSessionUsecase struct {
	date        usecase.DateTimeProvider
	sessionRepo usecase.SessionRepo
}

func NewAuthSession(
	date usecase.DateTimeProvider,
	sessionRepo usecase.SessionRepo,
) *AuthSessionUsecase {
	return &AuthSessionUsecase{
		date:        date,
		sessionRepo: sessionRepo,
	}
}

// Execute выполняет аутентификацию по сессии
func (u *AuthSessionUsecase) Execute(ctx usecase.AppCtx, sessionID string) (domain.AuthUser, error) {
	var zero domain.AuthUser

	// поиск сессии
	session, err := u.sessionRepo.Find(ctx, sessionID)
	switch err {
	case nil:
	case usecase.ErrNoData:
		return zero, usecase.ErrUnauthorized
	default:
		return zero, err
	}

	// сессия истекла
	if session.IsExpired(u.date.Now()) {
		return zero, usecase.ErrUnauthorized
	}

	return domain.NewAuthUser(session.UserID), nil
}
