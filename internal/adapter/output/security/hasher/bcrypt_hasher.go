package hasher

import (
	"event-driven-architecture/internal/usecase"

	"golang.org/x/crypto/bcrypt"
)

type bcryptHasher struct {
	cost int
}

func NewBCrypt(cost int) usecase.PasswordHasher {
	return &bcryptHasher{cost: cost}
}

func (h *bcryptHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (h *bcryptHasher) Compare(password, hash string) (bool, error) {
	switch err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err {
	case nil:
		return true, nil
	case bcrypt.ErrMismatchedHashAndPassword:
		return false, nil
	default:
		return false, err
	}
}
