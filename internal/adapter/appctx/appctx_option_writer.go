package appctx

import "event-driven-architecture/internal/usecase"

type appCtxOptionWriter struct {
	opts *options
}

func newAppCtxOptionWriter(opts *options) usecase.AppCtxOptionWriter {
	return &appCtxOptionWriter{opts: opts}
}

func (w *appCtxOptionWriter) SetSession(session usecase.Session) {
	w.opts.session = session
}
