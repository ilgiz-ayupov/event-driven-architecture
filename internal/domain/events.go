package domain

type UserCreatedEvent struct {
	Email  string
	userID string
}

func (e UserCreatedEvent) EventType() string {
	return "user.created"
}

func (e UserCreatedEvent) UserID() string {
	return e.userID
}

func NewUserCreatedEvent(email, userID string) UserCreatedEvent {
	return UserCreatedEvent{
		Email:  email,
		userID: userID,
	}
}
