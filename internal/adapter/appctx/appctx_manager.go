package appctx

import (
	"context"
	"event-driven-architecture/internal/usecase"
	"time"
)

type appCtxManager struct {
	timeout time.Duration
}

func NewAppCtxManager(timeout time.Duration) usecase.AppCtxManager {
	return &appCtxManager{timeout: timeout}
}

func (cm appCtxManager) CreateContext(parent context.Context, opts ...usecase.AppCtxOption) (usecase.AppCtx, context.CancelFunc) {
	options := options{}
	writer := newAppCtxOptionWriter(&options)

	usecase.ApplyAppCtxOptions(writer, opts...)

	ctx, cancel := context.WithTimeout(parent, cm.timeout)

	if options.session != nil {
		ctx = context.WithValue(ctx, SessionKey{}, options.session)
	}

	return newContext(ctx), cancel
}
