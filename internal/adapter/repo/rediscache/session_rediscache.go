package rediscache

import (
	"encoding/json"
	"event-driven-architecture/internal/domain"
	"event-driven-architecture/internal/usecase"
	"time"

	"github.com/redis/go-redis/v9"
)

type sessionRepo struct {
	client *redis.Client
}

func NewSession(client *redis.Client) usecase.SessionRepo {
	return &sessionRepo{client: client}
}

func (r *sessionRepo) Create(ctx usecase.AppCtx, session domain.Session) error {
	dto := newSessionDTO(
		session.ID,
		session.UserID,
		session.Roles,
		session.ExpiresAt,
	)

	data, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	ttl := time.Until(session.ExpiresAt)
	if ttl < 0 {
		// сессия уже истекла
		return nil
	}

	return r.client.Set(
		ctx.Context(),
		r.key(session.ID),
		data,
		ttl,
	).Err()
}

func (r *sessionRepo) Find(ctx usecase.AppCtx, sessionID string) (domain.Session, error) {
	var zero domain.Session

	data, err := r.client.Get(ctx.Context(), r.key(sessionID)).Result()
	switch err {
	case nil:
	case redis.Nil:
		return zero, usecase.ErrNoData
	default:
		return zero, err
	}

	var dto sessionDTO
	if err := json.Unmarshal([]byte(data), &dto); err != nil {
		return zero, err
	}

	return domain.NewSession(
		dto.ID,
		dto.UserID,
		dto.ExpiresAt,
	), nil
}

func (r *sessionRepo) Delete(ctx usecase.AppCtx, sessionID string) error {
	return r.client.Del(ctx.Context(), r.key(sessionID)).Err()
}

func (r *sessionRepo) key(sessionID string) string {
	return "sessions:" + sessionID
}
