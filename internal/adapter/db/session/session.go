package session

import (
	"event-driven-architecture/internal/usecase"

	"github.com/jmoiron/sqlx"
)

type session struct {
	conn    *sqlx.DB
	current *sqlx.Tx
}

func newSession(postgresConn *sqlx.DB) usecase.Session {
	return &session{conn: postgresConn}
}

func (s *session) Start() error {
	if s.current != nil {
		return nil
	}

	ts, err := s.conn.Beginx()
	if err != nil {
		return err
	}

	s.current = ts
	return nil
}

func (s *session) Commit() error {
	if s.current == nil {
		return nil
	}

	err := s.current.Commit()
	s.current = nil
	return err
}

func (s *session) Rollback() error {
	if s.current == nil {
		return nil
	}

	err := s.current.Rollback()
	s.current = nil
	return err
}

func (s *session) Tx() *sqlx.Tx {
	return s.current
}
