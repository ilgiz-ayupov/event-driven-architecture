package rediscache

import "time"

type sessionDTO struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	// TODO: реализовать роли
	Roles     []string  `json:"roles"`
	ExpiresAt time.Time `json:"expires_at"`
}

func newSessionDTO(id, userID string, roles []string, expiresAt time.Time) sessionDTO {
	return sessionDTO{
		ID:        id,
		UserID:    userID,
		Roles:     roles,
		ExpiresAt: expiresAt,
	}
}
