package usecase

func ApplyAppCtxOptions(
	writer AppCtxOptionWriter,
	opts ...AppCtxOption,
) {
	for _, opt := range opts {
		opt.apply(writer)
	}
}

type withSession struct {
	session Session
}

func (o withSession) apply(w AppCtxOptionWriter) {
	w.SetSession(o.session)
}

func WithSession(session Session) AppCtxOption {
	if session == nil {
		panic("session is nil")
	}
	return withSession{session: session}
}
