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
	SetTransaction(ts Transaction)
}

type AppCtx interface {
	Context() context.Context
	Transaction() Transaction
}

type AppCtxManager interface {
	CreateContext(parent context.Context, opts ...AppCtxOption) (AppCtx, context.CancelFunc)
}

type Transaction interface {
	Start() error
	Rollback() error
	Commit() error
	Tx() *sqlx.Tx
}

type TransactionManager interface {
	CreateTransaction() Transaction
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
	UserID() string
}

type EventPublisher interface {
	Publish(ctx AppCtx, event Event) error
}

type EventSubscription interface {
	Channel() <-chan []byte
	Close()
	ID() string
}

type EventBroker interface {
	Subscribe(sessionID string) EventSubscription
	SendToSession(sessionID string, msg []byte)
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

type SessionRepo interface {
	Create(ctx AppCtx, session domain.Session) error
	Find(ctx AppCtx, sessionID string) (domain.Session, error)
	Delete(ctx AppCtx, sessionID string) error
}

type SessionIndexRepo interface {
	Add(ctx AppCtx, sessionID, userID string) error
	Remove(ctx AppCtx, sessionID string) error
	SessionsByUser(ctx AppCtx, userID string) ([]string, error)
}
