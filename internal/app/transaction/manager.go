package transaction

import (
	"event-driven-architecture/internal/usecase"

	"github.com/jmoiron/sqlx"
)

type transactionManager struct {
	conn *sqlx.DB
}

func NewManager(postgresConn *sqlx.DB) usecase.TransactionManager {
	return &transactionManager{
		conn: postgresConn,
	}
}

func (sm *transactionManager) CreateTransaction() usecase.Transaction {
	return newTransaction(sm.conn)
}

func (sm *transactionManager) Close() error {
	return sm.conn.Close()
}
