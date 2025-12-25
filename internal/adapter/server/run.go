package server

import (
	"event-driven-architecture/internal/usecase"

	"github.com/gin-gonic/gin"
)

type ReturnFunc[T any] func(ctx usecase.AppCtx) (T, error)
type ExecFunc func(ctx usecase.AppCtx) error

func ReturnWithCommit[T any](
	c *gin.Context,
	transactionManager usecase.TransactionManager,
	appCtxManager usecase.AppCtxManager,
	fn ReturnFunc[T],
) (T, error) {
	var zero T

	// открыть сессию
	ts := transactionManager.CreateTransaction()
	if err := ts.Start(); err != nil {
		return zero, usecase.ErrInternalError
	}
	defer ts.Rollback()

	// сформировать контекст
	appCtx, cancel := appCtxManager.CreateContext(
		c.Request.Context(),
		usecase.WithTransaction(ts),
	)
	defer cancel()

	// выполнить бизнес-логику
	data, err := fn(appCtx)
	if err != nil {
		return zero, err
	}

	// зафиксировать транзакцию
	if err := ts.Commit(); err != nil {
		return zero, usecase.ErrInternalError
	}

	return data, nil
}

func Exec(
	c *gin.Context,
	transactionManager usecase.TransactionManager,
	appCtxManager usecase.AppCtxManager,
	fn ExecFunc,
) error {
	// открыть сессию
	ts := transactionManager.CreateTransaction()
	if err := ts.Start(); err != nil {
		return usecase.ErrInternalError
	}
	defer ts.Rollback()

	// сформировать контекст
	appCtx, cancel := appCtxManager.CreateContext(
		c.Request.Context(),
		usecase.WithTransaction(ts),
	)
	defer cancel()

	// выполнить бизнес-логику
	err := fn(appCtx)
	if err != nil {
		return err
	}

	// зафиксировать транзакцию
	if err := ts.Commit(); err != nil {
		return usecase.ErrInternalError
	}

	return nil

}
