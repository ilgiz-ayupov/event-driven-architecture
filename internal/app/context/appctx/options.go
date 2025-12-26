package appctx

import "event-driven-architecture/internal/usecase"

type options struct {
	ts usecase.Transaction
}

type appCtxOptionWriter struct {
	opts *options
}

func newAppCtxOptionWriter(opts *options) usecase.AppCtxOptionWriter {
	return &appCtxOptionWriter{opts: opts}
}

func (w *appCtxOptionWriter) SetTransaction(ts usecase.Transaction) {
	w.opts.ts = ts
}
