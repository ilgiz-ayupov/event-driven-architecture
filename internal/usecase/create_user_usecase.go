package usecase

import (
	"context"
	"event-driven-architecture/internal/domain"
)

type CreateUserUseCase struct {
	log            Logger
	sessionManager SessionManager
	appCtxManager  AppCtxManager
	eventPublisher EventPublisher
	idGenerator    IDGenerator
	passwordHasher PasswordHasher

	userRepo UserRepo
}

func NewCreateUser(
	log Logger,
	sessionManager SessionManager,
	appCtxManager AppCtxManager,
	eventPublisher EventPublisher,
	idGenerator IDGenerator,
	passwordHasher PasswordHasher,

	userRepo UserRepo,
) *CreateUserUseCase {
	return &CreateUserUseCase{
		log:            log,
		sessionManager: sessionManager,
		appCtxManager:  appCtxManager,
		eventPublisher: eventPublisher,
		idGenerator:    idGenerator,
		passwordHasher: passwordHasher,

		userRepo: userRepo,
	}
}

func (u *CreateUserUseCase) Execute(ctx context.Context, email, password string) error {
	// открыть сессию
	session := u.sessionManager.CreateSession()
	if err := session.Start(); err != nil {
		u.log.Error("не удалось открыть транзакцию", "error", err)
		return ErrInternalError
	}
	defer session.Rollback()

	// сформировать контекст
	appCtx, cancel := u.appCtxManager.CreateContext(
		ctx,
		WithSession(session),
	)
	defer cancel()

	// хешировать пароль
	hash, err := u.passwordHasher.Hash(password)
	if err != nil {
		u.log.Error("не удалось хешировать пароль", "error", err)
		return ErrInternalError
	}

	// сохранить пользователя
	user := domain.NewUser(
		u.idGenerator.NewID(),
		email,
		hash,
	)

	if err := u.userRepo.Create(appCtx, user); err != nil {
		u.log.Error("не удалось сохранить пользователя", "error", err)
		return ErrInternalError
	}

	// опубликовать событие - пользователь создан
	event := domain.NewUserCreatedEvent(
		user.ID,
		user.Email,
	)

	if err := u.eventPublisher.Publish(event); err != nil {
		u.log.Error(
			"не удалось опубликовать событие",
			"error", err,
			"event_type", event.EventType(),
		)

		return ErrInternalError
	} else {
		u.log.Info(
			"событие опубликовано",
			"event_type", event.EventType(),
		)
	}

	// зафиксировать сессию
	if err := session.Commit(); err != nil {
		u.log.Error("не удалось зафиксировать сессию", "error", err)
		return ErrInternalError
	}

	// вернуть результат
	u.log.Info(
		"пользователь создан",
		"user_id", user.ID,
	)

	return nil
}
