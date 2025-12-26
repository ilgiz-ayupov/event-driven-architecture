package postgres

import (
	"event-driven-architecture/internal/usecase"

	"github.com/jmoiron/sqlx"
)

func SqlxTx(transaction usecase.Transaction) *sqlx.Tx {
	return transaction.Tx()
}
