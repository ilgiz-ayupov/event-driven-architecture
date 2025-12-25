package appctx

import "event-driven-architecture/internal/usecase"

type options struct {
	ts usecase.Transaction
}

type TransactionKey struct{}
