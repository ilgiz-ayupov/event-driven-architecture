package auth_usecase

import (
	"event-driven-architecture/internal/domain"
	"event-driven-architecture/internal/usecase"
	"time"
)

type LoginUserUseCase struct {
	log            usecase.Logger
	idGenerator    usecase.IDGenerator
	passwordHasher usecase.PasswordHasher

	userRepo         usecase.UserRepo
	sessionRepo      usecase.SessionRepo
	sessionIndexRepo usecase.SessionIndexRepo
}

func NewLoginUser(
	log usecase.Logger,
	idGenerator usecase.IDGenerator,
	passwordHasher usecase.PasswordHasher,

	userRepo usecase.UserRepo,
	sessionRepo usecase.SessionRepo,
	sessionIndexRepo usecase.SessionIndexRepo,
) *LoginUserUseCase {
	return &LoginUserUseCase{
		log:            log,
		idGenerator:    idGenerator,
		passwordHasher: passwordHasher,

		userRepo:         userRepo,
		sessionRepo:      sessionRepo,
		sessionIndexRepo: sessionIndexRepo,
	}
}

func (u *LoginUserUseCase) Execute(ctx usecase.AppCtx, email, password string) (sessionID string, err error) {
	var zero string

	// создать сессию
	sessionID = u.idGenerator.NewID()

	session := domain.Session{
		ID:        sessionID,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := u.sessionRepo.Create(ctx, session); err != nil {
		u.log.Error("не удалось создать сессию", "error", err)
		return zero, usecase.ErrInternalError
	}

	if err := u.sessionIndexRepo.Add(ctx, sessionID, user.ID); err != nil {
		u.log.Error("не удалось создать индекс сессии", "error", err)
		return zero, usecase.ErrInternalError
	}

	return sessionID, nil
}
