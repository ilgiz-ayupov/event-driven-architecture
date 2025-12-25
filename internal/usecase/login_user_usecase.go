package usecase

import (
	"event-driven-architecture/internal/domain"
	"time"
)

type LoginUserUseCase struct {
	log            Logger
	idGenerator    IDGenerator
	passwordHasher PasswordHasher

	userRepo         UserRepo
	sessionRepo      SessionRepo
	sessionIndexRepo SessionIndexRepo
}

func NewLoginUser(
	log Logger,
	idGenerator IDGenerator,
	passwordHasher PasswordHasher,

	userRepo UserRepo,
	sessionRepo SessionRepo,
	sessionIndexRepo SessionIndexRepo,
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

func (u *LoginUserUseCase) Execute(ctx AppCtx, email, password string) (sessionID string, err error) {
	var zero string

	// поиск пользователя
	user, err := u.userRepo.FindByEmail(ctx, email)
	switch err {
	case nil:
	case ErrNoData:
		u.log.Warn("пользователь с таким email не найден", "email", email, "error", err)
		return zero, ErrInvalidCredentials
	default:
		u.log.Error("не удалось найти пользователя по email", "email", email, "error", err)
		return zero, ErrInternalError
	}

	// проверка пароля
	switch ok, err := u.passwordHasher.Compare(password, user.PasswordHash); err {
	case nil:
		if !ok {
			u.log.Warn("неверный пароль при авторизации пользователя", "email", email)
			return zero, ErrInvalidCredentials
		}
	default:
		u.log.Error("не удалось проверить пароль при авторизации пользователя", "error", err)
		return zero, ErrInternalError
	}

	// создать сессию
	sessionID = u.idGenerator.NewID()

	session := domain.Session{
		ID:        sessionID,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := u.sessionRepo.Create(ctx, session); err != nil {
		u.log.Error("не удалось создать сессию", "error", err)
		return zero, ErrInternalError
	}

	if err := u.sessionIndexRepo.Add(ctx, sessionID, user.ID); err != nil {
		u.log.Error("не удалось создать индекс сессии", "error", err)
		return zero, ErrInternalError
	}

	return sessionID, nil
}
