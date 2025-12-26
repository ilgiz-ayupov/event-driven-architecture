package user_usecase

type CreateUserInput struct {
	Email    string
	Password string
}

func NewCreateUserInput(
	email string,
	password string,
) CreateUserInput {
	return CreateUserInput{
		Email:    email,
		Password: password,
	}
}

type CreateUserOutput struct {
	UserID string
}

func NewCreateUserOutput(userID string) CreateUserOutput {
	return CreateUserOutput{
		UserID: userID,
	}
}
