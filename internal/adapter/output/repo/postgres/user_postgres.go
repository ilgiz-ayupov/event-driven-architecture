package postgres

import (
	"database/sql"
	"event-driven-architecture/internal/domain"
	"event-driven-architecture/internal/usecase"
)

type userRepo struct {
	log usecase.Logger
}

func NewUser(log usecase.Logger) usecase.UserRepo {
	return &userRepo{log: log}
}

// Create создает пользователя
func (r *userRepo) Create(ctx usecase.AppCtx, user domain.User) error {
	q := `
		INSERT INTO users (user_id, email, password_hash, created_at)
		VALUES (:user_id, :email, :password_hash, NOW())
	`

	_, err := SqlxTx(ctx.Transaction()).NamedExecContext(ctx.Context(), q, map[string]any{
		"user_id":       user.ID,
		"email":         user.Email,
		"password_hash": user.PasswordHash,
	})
	if err != nil {
		r.log.Error(
			"не удалось создать пользователя в Postgres",
			"user_id", user.ID,
			"email", user.Email,
			"error", err,
		)

		return usecase.ErrInternalError
	}

	return nil
}

// FindByEmail ищет пользователя по email
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
	err := SqlxTx(ctx.Transaction()).QueryRowxContext(ctx.Context(), q, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
	)
	switch err {
	case nil:
		return user, err
	case sql.ErrNoRows:
		r.log.Warn(
			"пользователь не найден по email в Postgres",
			"email", email,
		)
		return zero, usecase.ErrNoData
	default:
		r.log.Error(
			"не удалось найти пользователя по email в Postgres",
			"email", email,
			"error", err,
		)
		return zero, err
	}
}
