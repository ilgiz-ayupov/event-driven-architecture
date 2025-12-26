package auth_usecase

type AuthUserInput struct {
	Email    string
	Password string
}

func NewAuthUserInput(email, password string) AuthUserInput {
	return AuthUserInput{
		Email:    email,
		Password: password,
	}
}
