package postgres

import (
	"event-driven-architecture/internal/domain"
	"event-driven-architecture/internal/usecase"
	"fmt"
)

type userRepo struct{}

func NewUser() usecase.UserRepo {
	return &userRepo{}
}

func (r *userRepo) Create(ctx usecase.AppCtx, user domain.User) error {
	q := `
		INSERT INTO users (user_id, email, password_hash)
		VALUES (:user_id, :email, :password_hash)
	`

	_, err := ctx.Session().Tx().NamedExecContext(ctx.Context(), q, map[string]any{
		"user_id":       user.ID,
		"email":         user.Email,
		"password_hash": user.PasswordHash,
	})
	return err
}

func (r *userRepo) FindByEmail(ctx usecase.AppCtx, email string) (domain.User, error) {
	var zero domain.User
	return zero, fmt.Errorf("not implemented")
}
