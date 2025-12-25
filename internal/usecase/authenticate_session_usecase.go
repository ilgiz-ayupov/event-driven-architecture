package usecase

import (
	"event-driven-architecture/internal/domain"
	"time"
)

type AuthenticateSessionUseCase struct {
	log         Logger
	sessionRepo SessionRepo
}

func NewAuthenticateSession(
	log Logger,
	sessionRepo SessionRepo,
) *AuthenticateSessionUseCase {
	return &AuthenticateSessionUseCase{
		log:         log,
		sessionRepo: sessionRepo,
	}
}

func (u *AuthenticateSessionUseCase) Execute(ctx AppCtx, sessionID string) (domain.AuthUser, error) {
	var zero domain.AuthUser

	// поиск сессии
	session, err := u.sessionRepo.Find(ctx, sessionID)
	switch err {
	case nil:
	case ErrNoData:
		return zero, ErrUnauthorized
	default:
		u.log.Error("не удалось найти сессию в кеше", "error", err)
		return zero, ErrInternalError
	}

	// сессия истекла
	// TODO: вывести time.Now в адаптер
	if session.IsExpired(time.Now()) {
		return zero, ErrUnauthorized
	}

	return domain.NewAuthUser(session.UserID), nil
}
