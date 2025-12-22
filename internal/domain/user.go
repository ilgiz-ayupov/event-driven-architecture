package domain

type User struct {
	ID           string
	Email        string
	PasswordHash string
}

func NewUser(id, email, passwordHash string) User {
	return User{
		ID:           id,
		Email:        email,
		PasswordHash: passwordHash,
	}
}
