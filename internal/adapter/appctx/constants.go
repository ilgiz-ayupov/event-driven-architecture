package appctx

import "event-driven-architecture/internal/usecase"

type options struct {
	session usecase.Session
}

type SessionKey struct{}
