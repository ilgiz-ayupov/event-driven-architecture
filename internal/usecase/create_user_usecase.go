package usecase

import (
	"event-driven-architecture/internal/domain"
)

type CreateUserUseCase struct {
	log            Logger
	eventPublisher EventPublisher
	idGenerator    IDGenerator
	passwordHasher PasswordHasher

	userRepo UserRepo
}

func NewCreateUser(
	log Logger,
	eventPublisher EventPublisher,
	idGenerator IDGenerator,
	passwordHasher PasswordHasher,

	userRepo UserRepo,
) *CreateUserUseCase {
	return &CreateUserUseCase{
		log:            log,
		eventPublisher: eventPublisher,
		idGenerator:    idGenerator,
		passwordHasher: passwordHasher,

		userRepo: userRepo,
	}
}

type CreateUserInput struct {
	Email    string
	Password string

	UserLog UserLog
}

func NewCreateUserInput(
	email string,
	password string,
	userLog UserLog,
) CreateUserInput {
	return CreateUserInput{
		Email:    email,
		Password: password,

		UserLog: userLog,
	}
}

func (u *CreateUserUseCase) Execute(ctx AppCtx, in CreateUserInput) error {
	// хешировать пароль
	hash, err := u.passwordHasher.Hash(in.Password)
	if err != nil {
		u.log.Error("не удалось хешировать пароль", "error", err)
		return ErrInternalError
	}

	// сохранить пользователя
	user := domain.NewUser(
		u.idGenerator.NewID(),
		in.Email,
		hash,
	)

	if err := u.userRepo.Create(ctx, user); err != nil {
		u.log.Error("не удалось сохранить пользователя", "error", err)
		return ErrInternalError
	}

	// опубликовать событие - пользователь создан
	event := domain.NewUserCreatedEvent(
		in.Email,
		in.UserLog.UserID,
	)

	if err := u.eventPublisher.Publish(ctx, event); err != nil {
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

	// вернуть результат
	u.log.Info(
		"пользователь создан",
		"user_id", in.UserLog.UserID,
	)

	return nil
}
