package postgres

import (
	"event-driven-architecture/internal/domain"
	"event-driven-architecture/internal/usecase"
)

type sessionRepo struct{}

func NewSession() usecase.SessionRepo {
	return &sessionRepo{}
}

func (r *sessionRepo) Create(ctx usecase.AppCtx, session domain.Session) error {
	q := `
		INSERT INTO sessions (session_id, user_id, created_at, expires_at)
		VALUES (:session_id, :user_id, NOW(), :expires_at)
	`

	_, err := ctx.Transaction().Tx().NamedExecContext(ctx.Context(), q, map[string]any{
		"session_id": session.ID,
		"user_id":    session.UserID,
		"expires_at": session.ExpiresAt,
	})

	return err
}

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
	err := ctx.Transaction().Tx().QueryRowxContext(ctx.Context(), q, sessionID).Scan(
		&session.ID,
		&session.UserID,
		&session.ExpiresAt,
	)

	return session, err
}

func (r *sessionRepo) Delete(ctx usecase.AppCtx, sessionID string) error {
	q := `
		DELETE FROM sessions
		WHERE session_id = :session_id
	`

	_, err := ctx.Transaction().Tx().NamedExecContext(ctx.Context(), q, map[string]any{
		"session_id": sessionID,
	})

	return err
}
