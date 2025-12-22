package session

import (
	"event-driven-architecture/internal/usecase"

	"github.com/jmoiron/sqlx"
)

type sessionManager struct {
	conn *sqlx.DB
}

func NewSessionManager(postgresConn *sqlx.DB) usecase.SessionManager {
	return &sessionManager{
		conn: postgresConn,
	}
}

func (sm *sessionManager) CreateSession() usecase.Session {
	return newSession(sm.conn)
}

func (sm *sessionManager) Close() error {
	return sm.conn.Close()
}
