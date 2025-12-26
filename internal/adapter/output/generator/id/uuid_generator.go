package id

import (
	"event-driven-architecture/internal/usecase"

	"github.com/google/uuid"
)

type uuidGenerator struct{}

func NewUUIDGenerator() usecase.IDGenerator {
	return &uuidGenerator{}
}

func (g *uuidGenerator) NewID() string {
	return uuid.NewString()
}
