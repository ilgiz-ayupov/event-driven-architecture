package postgres

import (
	"database/sql"
	"event-driven-architecture/internal/domain"
	"event-driven-architecture/internal/usecase"
)

type userRepo struct{}

func NewUser() usecase.UserRepo {
	return &userRepo{}
}

func (r *userRepo) Create(ctx usecase.AppCtx, user domain.User) error {
	q := `
		INSERT INTO users (user_id, email, password_hash, created_at)
		VALUES (:user_id, :email, :password_hash, NOW())
	`

	_, err := ctx.Transaction().Tx().NamedExecContext(ctx.Context(), q, map[string]any{
		"user_id":       user.ID,
		"email":         user.Email,
		"password_hash": user.PasswordHash,
	})
	return err
}

func (r *userRepo) FindByEmail(ctx usecase.AppCtx, email string) (domain.User, error) {
	var zero domain.User

	q := `
		SELECT
			u.user_id
			, u.email
			, u.password_hash
		FROM users u
		WHERE u.email = $1
	`

	var user domain.User
	err := ctx.Transaction().Tx().QueryRowxContext(ctx.Context(), q, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
	)
	switch err {
	case nil:
		return user, err
	case sql.ErrNoRows:
		return zero, usecase.ErrNoData
	default:
		return zero, err
	}

}
