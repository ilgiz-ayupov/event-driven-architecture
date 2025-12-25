package usecase

func ApplyAppCtxOptions(
	writer AppCtxOptionWriter,
	opts ...AppCtxOption,
) {
	for _, opt := range opts {
		opt.apply(writer)
	}
}

type withTransaction struct {
	ts Transaction
}

func (o withTransaction) apply(w AppCtxOptionWriter) {
	w.SetTransaction(o.ts)
}

func WithTransaction(ts Transaction) AppCtxOption {
	if ts == nil {
		panic("transaction is nil")
	}
	return withTransaction{ts: ts}
}
