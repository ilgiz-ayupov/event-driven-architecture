package domain

type AuthUser struct {
	ID string
	// TODO: реализовать роли в будущем
	Roles []string
}

func NewAuthUser(id string) AuthUser {
	return AuthUser{
		ID: id,
	}
}
