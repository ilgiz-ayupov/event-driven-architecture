package broker

import (
	"event-driven-architecture/internal/usecase"
	"sync"

	"github.com/google/uuid"
)

type subscriber struct {
	ch   chan []byte
	once sync.Once
}

type sseBroker struct {
	log usecase.Logger

	mu          sync.RWMutex
	subscribers map[string]*subscriber
}

func NewSSEBroker(log usecase.Logger) usecase.EventBroker {
	return &sseBroker{
		log:         log,
		subscribers: make(map[string]*subscriber),
	}
}

func (b *sseBroker) Subscribe() (id string, subscriberCh <-chan []byte) {
	b.mu.Lock()

	id = uuid.NewString()
	sub := &subscriber{ch: make(chan []byte, 16)}
	b.subscribers[id] = sub

	b.mu.Unlock()

	b.log.Info(
		"SSE-брокер: подписка клиента",
		"subscriber_id", id,
	)
	return id, sub.ch
}

func (b *sseBroker) Unsubscribe(id string) {
	b.mu.Lock()

	sub, ok := b.subscribers[id]
	if !ok {
		b.mu.Unlock()
		return
	}

	delete(b.subscribers, id)
	sub.once.Do(func() { close(sub.ch) })

	b.log.Info(
		"SSE-брокер: отписка клиента",
		"subscriber_id", id,
	)
	b.mu.Unlock()
}

func (b *sseBroker) Broadcast(msg []byte) {
	// snapshot подписчиков
	b.mu.RLock()

	subs := make([]struct {
		id  string
		sub *subscriber
	}, 0, len(b.subscribers))

	for id, sub := range b.subscribers {
		subs = append(subs, struct {
			id  string
			sub *subscriber
		}{id: id, sub: sub})
	}

	b.mu.RUnlock()

	// рассылка
	var toDrop []string

	for _, s := range subs {
		select {
		case s.sub.ch <- msg:
			// ok
			b.log.Info(
				"SSE-брокер: отправка сообщения клиенту",
				"subscriber_id", s.id,
			)
		default:
			// клиент не читает - помечаем на дроп
			b.log.Info(
				"SSE-брокер: клиент помечен на дроп",
				"subscriber_id", s.id,
			)
			toDrop = append(toDrop, s.id)
		}
	}

	// дропаем помеченных клиентов
	if len(toDrop) > 0 {
		b.mu.Lock()

		for _, id := range toDrop {
			sub, ok := b.subscribers[id]
			if ok {
				delete(b.subscribers, id)
				sub.once.Do(func() { close(sub.ch) })

				b.log.Info(
					"SSE-брокер: клиент дропнут",
					"subscriber_id", id,
				)
			}
		}

		b.mu.Unlock()
	}
}
