package domain

import "time"

type Session struct {
	ID     string
	UserID string
	// TODO: реализовать рол
	Roles     []string
	ExpiresAt time.Time
}

func (s Session) IsExpired(now time.Time) bool {
	return now.After(s.ExpiresAt)
}

func NewSession(id, userID string, expiresAt time.Time) Session {
	return Session{
		ID:        id,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}
}
