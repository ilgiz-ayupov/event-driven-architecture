package broker

import (
	"event-driven-architecture/internal/usecase"
	"sync"

	"github.com/google/uuid"
)

type sseBroker struct {
	log usecase.Logger

	mu sync.RWMutex
	// session_id -> set of subscribers
	sessions map[string]map[*connection]struct{}
}

func NewSSEBroker(log usecase.Logger) usecase.EventBroker {
	return &sseBroker{
		log:      log,
		sessions: make(map[string]map[*connection]struct{}),
	}
}

func (b *sseBroker) Subscribe(sessionID string) usecase.EventSubscription {
	subscriptionID := uuid.NewString()
	conn := newConnection(make(chan []byte, 16))

	b.mu.Lock()
	if _, ok := b.sessions[sessionID]; !ok {
		b.sessions[sessionID] = make(map[*connection]struct{})
	}
	b.sessions[sessionID][conn] = struct{}{}
	b.mu.Unlock()

	b.log.Info(
		"SSE-брокер: подписка",
		"session_id", sessionID,
		"subscription_id", subscriptionID,
		"total_subscribers", len(b.sessions[sessionID]),
	)

	return newSubscription(
		subscriptionID,
		sessionID,
		conn,
		b,
	)
}

func (b *sseBroker) unsubscribe(
	sessionID string,
	sub *connection,
	subscriptionID string,
) {
	b.mu.Lock()

	subs, ok := b.sessions[sessionID]
	if !ok {
		b.mu.Unlock()
		return
	}

	if _, ok := subs[sub]; !ok {
		b.mu.Unlock()
		return
	}

	delete(subs, sub)
	sub.once.Do(func() {
		close(sub.ch)
	})

	if len(subs) == 0 {
		delete(b.sessions, sessionID)
	}

	b.log.Info(
		"SSE-брокер: отписка",
		"session_id", sessionID,
		"subscription_id", subscriptionID,
		"total_subscribers", len(b.sessions[sessionID]),
	)
	b.mu.Unlock()
}

func (b *sseBroker) SendToSession(sessionID string, msg []byte) {
	// snapshot подписок сессий
	b.mu.RLock()
	subs, ok := b.sessions[sessionID]
	if !ok {
		b.mu.RUnlock()
		return
	}

	targets := make([]*connection, 0, len(subs))
	for instance := range subs {
		targets = append(targets, instance)
	}
	b.mu.RUnlock()

	b.log.Info(
		"SSE-брокер: количество подписчиков для сессии",
		"session_id", sessionID,
		"total_subscribers", len(targets),
	)

	// рассылка
	var toDrop []*connection

	for _, s := range targets {
		select {
		case s.ch <- msg:
			// ok
		default:
			// клиент не читает - помечаем на дроп
			toDrop = append(toDrop, s)
		}
	}

	// дропаем помеченных клиентов
	if len(toDrop) > 0 {
		b.mu.Lock()
		for _, instance := range toDrop {
			if subscriptions, ok := b.sessions[sessionID]; ok {
				delete(subscriptions, instance)
				instance.once.Do(func() { close(instance.ch) })
			}
		}
		if subscriptions, ok := b.sessions[sessionID]; ok && len(subscriptions) == 0 {
			delete(b.sessions, sessionID)
		}
		b.mu.Unlock()
	}
}
