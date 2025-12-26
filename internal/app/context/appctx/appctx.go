package appctx

import (
	"context"
	"event-driven-architecture/internal/usecase"
)

type appCtx struct {
	ctx context.Context
}

func newContext(ctx context.Context) usecase.AppCtx {
	return &appCtx{ctx: ctx}
}

func (c appCtx) Context() context.Context {
	return c.ctx
}

func (c appCtx) Transaction() usecase.Transaction {
	ts, ok := c.ctx.Value(TransactionKey{}).(usecase.Transaction)
	if !ok {
		panic(errTransactionNotFound)
	}

	return ts
}
