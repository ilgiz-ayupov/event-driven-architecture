package usecase

import (
	"context"
	"event-driven-architecture/internal/domain"

	"github.com/jmoiron/sqlx"
)

type AppCtxOption interface {
	apply(AppCtxOptionWriter)
}

type AppCtxOptionWriter interface {
	SetSession(session Session)
}

type AppCtx interface {
	Context() context.Context
	Session() Session
}

type AppCtxManager interface {
	CreateContext(parent context.Context, opts ...AppCtxOption) (AppCtx, context.CancelFunc)
}

type Session interface {
	Start() error
	Rollback() error
	Commit() error
	Tx() *sqlx.Tx
}

type SessionManager interface {
	CreateSession() Session
	Close() error
}

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type Event interface {
	EventType() string
}

type EventPublisher interface {
	Publish(event Event) error
}

type EventBroker interface {
	Subscribe() (subscriberID string, ch <-chan []byte)
	Unsubscribe(subscriberID string)
	Broadcast(msg []byte)
}

type IDGenerator interface {
	NewID() string
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(password, hash string) (bool, error)
}

type UserRepo interface {
	Create(ctx AppCtx, user domain.User) error
	FindByEmail(ctx AppCtx, email string) (domain.User, error)
}
