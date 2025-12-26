package appctx

import "errors"

var (
	errTransactionNotFound = errors.New("транзакция не найдена в контексте")
)
