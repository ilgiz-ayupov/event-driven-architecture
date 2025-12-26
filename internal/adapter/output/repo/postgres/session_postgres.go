package postgres

import (
	"event-driven-architecture/internal/domain"
	"event-driven-architecture/internal/usecase"
)

type sessionRepo struct {
	log usecase.Logger
}

func NewSession(log usecase.Logger) usecase.SessionRepo {
	return &sessionRepo{log: log}
}

// Create создает сессию
func (r *sessionRepo) Create(ctx usecase.AppCtx, session domain.Session) error {
	q := `
		INSERT INTO sessions (session_id, user_id, created_at, expires_at)
		VALUES (:session_id, :user_id, NOW(), :expires_at)
	`

	_, err := SqlxTx(ctx.Transaction()).NamedExecContext(ctx.Context(), q, map[string]any{
		"session_id": session.ID,
		"user_id":    session.UserID,
		"expires_at": session.ExpiresAt,
	})
	if err != nil {
		r.log.Error(
			"не удалось создать сессию в Postgres",
			"session_id", session.ID,
			"user_id", session.UserID,
			"error", err,
		)

		return usecase.ErrInternalError
	}

	return nil
}

// Find ищет сессию
func (r *sessionRepo) Find(ctx usecase.AppCtx, sessionID string) (domain.Session, error) {
	q := `
		SELECT
			s.session_id
			, s.user_id
			, s.expires_at
		FROM sessions s
		WHERE s.id = $1
			AND s.expires_at > NOW()
	`

	var session domain.Session
	err := SqlxTx(ctx.Transaction()).QueryRowxContext(ctx.Context(), q, sessionID).Scan(
		&session.ID,
		&session.UserID,
		&session.ExpiresAt,
	)
	if err != nil {
		r.log.Error(
			"не удалось найти сессию в Postgres",
			"session_id", sessionID,
			"error", err,
		)

		return session, usecase.ErrInternalError
	}

	return session, nil
}

// Delete удаляет сессию
func (r *sessionRepo) Delete(ctx usecase.AppCtx, sessionID string) error {
	q := `
		DELETE FROM sessions
		WHERE session_id = :session_id
	`

	_, err := SqlxTx(ctx.Transaction()).NamedExecContext(ctx.Context(), q, map[string]any{
		"session_id": sessionID,
	})
	if err != nil {
		r.log.Error(
			"не удалось удалить сессию из Postgres",
			"session_id", sessionID,
			"error", err,
		)

		return usecase.ErrInternalError
	}

	return nil
}
