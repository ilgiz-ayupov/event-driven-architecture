package publisher

import (
	"encoding/json"
	"event-driven-architecture/internal/usecase"
	"fmt"
)

type ssePublisher struct {
	broker           usecase.EventBroker
	sessionIndexRepo usecase.SessionIndexRepo
}

func NewSSEPublisher(
	broker usecase.EventBroker,
	sessionIndexRepo usecase.SessionIndexRepo,
) usecase.EventPublisher {
	return &ssePublisher{
		broker:           broker,
		sessionIndexRepo: sessionIndexRepo,
	}
}

func (p *ssePublisher) Publish(ctx usecase.AppCtx, event usecase.Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	userID := event.UserID()
	eventType := event.EventType()

	msg := makeSSEMessage(eventType, data)

	sessions, err := p.sessionIndexRepo.SessionsByUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get sessions for user %s: %w", userID, err)
	}

	for _, sessionID := range sessions {
		p.broker.SendToSession(sessionID, msg)
	}

	return nil
}

func makeSSEMessage(eventType string, data []byte) []byte {
	capacity := len("event: \ndata: \n\n") + len(eventType) + len(data)
	msg := make([]byte, 0, capacity)

	msg = append(msg, "event: "...)
	msg = append(msg, eventType...)
	msg = append(msg, "\ndata: "...)
	msg = append(msg, data...)

	return msg
}
