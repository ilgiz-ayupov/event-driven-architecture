package transaction

import (
	"event-driven-architecture/internal/usecase"

	"github.com/jmoiron/sqlx"
)

type transaction struct {
	conn    *sqlx.DB
	current *sqlx.Tx
}

func newTransaction(postgresConn *sqlx.DB) usecase.Transaction {
	return &transaction{conn: postgresConn}
}

func (s *transaction) Start() error {
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

func (s *transaction) Commit() error {
	if s.current == nil {
		return nil
	}

	err := s.current.Commit()
	s.current = nil
	return err
}

func (s *transaction) Rollback() error {
	if s.current == nil {
		return nil
	}

	err := s.current.Rollback()
	s.current = nil
	return err
}

func (s *transaction) Tx() *sqlx.Tx {
	return s.current
}
