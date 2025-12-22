package publisher

import (
	"event-driven-architecture/internal/usecase"
	"fmt"
)

type multiPublisher struct {
	publishers []usecase.EventPublisher
}

func NewMultiPublisher(publishers ...usecase.EventPublisher) usecase.EventPublisher {
	return &multiPublisher{
		publishers: publishers,
	}
}

func (p *multiPublisher) Publish(event usecase.Event) error {
	var errs []error

	for _, pub := range p.publishers {
		if err := pub.Publish(event); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("multi-publish errors: %v", errs)
	}

	return nil
}
