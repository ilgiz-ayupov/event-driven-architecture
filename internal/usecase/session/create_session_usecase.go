package session_usecase

import "event-driven-architecture/internal/usecase"

type CreateSessionUseCase struct {
	sessionRepo usecase.SessionRepo
	idGenerator usecase.IDGenerator
}

func NewCreateSession(
	sessionRepo usecase.SessionRepo,
	idGenerator usecase.IDGenerator,
) *CreateSessionUseCase {
	return &CreateSessionUseCase{
		sessionRepo: sessionRepo,
		idGenerator: idGenerator,
	}
}

func (u *CreateSessionUseCase) Execute(ctx usecase.AppCtx) error
