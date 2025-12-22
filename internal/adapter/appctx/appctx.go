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

func (c appCtx) Session() usecase.Session {
	session, ok := c.ctx.Value(SessionKey{}).(usecase.Session)
	if !ok {
		panic(errSessionNotFound)
	}

	return session
}
