package domain

type UserCreatedEvent struct {
	UserID string
	Email  string
}

func (UserCreatedEvent) EventType() string {
	return "user.created"
}

func NewUserCreatedEvent(userID, email string) UserCreatedEvent {
	return UserCreatedEvent{
		UserID: userID,
		Email:  email,
	}
}
