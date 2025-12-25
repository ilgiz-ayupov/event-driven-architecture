package broker

import (
	"event-driven-architecture/internal/usecase"
	"sync"
)

type subscription struct {
	id        string
	sessionID string
	conn      *connection
	broker    *sseBroker

	once sync.Once
}

func newSubscription(
	id string,
	sessionID string,
	conn *connection,
	broker *sseBroker,
) usecase.EventSubscription {
	return &subscription{
		id:        id,
		sessionID: sessionID,
		conn:      conn,
		broker:    broker,
	}
}

func (s *subscription) ID() string {
	return s.id
}

func (s *subscription) Channel() <-chan []byte {
	return s.conn.ch
}

func (s *subscription) Close() {
	s.once.Do(func() {
		s.broker.unsubscribe(s.sessionID, s.conn, s.id)
	})
}
