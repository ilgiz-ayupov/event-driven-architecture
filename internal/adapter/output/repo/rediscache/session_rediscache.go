package rediscache

import (
	"encoding/json"
	"event-driven-architecture/internal/domain"
	"event-driven-architecture/internal/usecase"
	"time"

	"github.com/redis/go-redis/v9"
)

type sessionRepo struct {
	log    usecase.Logger
	client *redis.Client
}

func NewSession(log usecase.Logger, client *redis.Client) usecase.SessionRepo {
	return &sessionRepo{log: log, client: client}
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
		r.log.Error(
			"не удалось сериализовать сессию",
			"session_id", session.ID,
			"error", err,
		)
		return usecase.ErrInternalError
	}

	ttl := time.Until(session.ExpiresAt)
	if ttl < 0 {
		r.log.Warn(
			"сессия уже истекла, не сохраняем в кэш",
			"session_id", session.ID,
		)
		return nil
	}

	if err := r.client.Set(
		ctx.Context(),
		r.key(session.ID),
		data,
		ttl,
	).Err(); err != nil {
		r.log.Error(
			"не удалось сохранить сессию в Redis",
			"session_id", session.ID,
			"error", err,
		)
		return usecase.ErrInternalError
	}

	return nil
}

func (r *sessionRepo) Find(ctx usecase.AppCtx, sessionID string) (domain.Session, error) {
	var zero domain.Session

	data, err := r.client.Get(ctx.Context(), r.key(sessionID)).Result()
	switch err {
	case nil:
	case redis.Nil:
		r.log.Warn(
			"сессия не найдена в кэше",
			"session_id", sessionID,
		)
		return zero, usecase.ErrNoData
	default:
		r.log.Error(
			"не удалось получить сессию из Redis",
			"session_id", sessionID,
			"error", err,
		)
		return zero, usecase.ErrInternalError
	}

	var dto sessionDTO
	if err := json.Unmarshal([]byte(data), &dto); err != nil {
		r.log.Error(
			"не удалось десериализовать сессию",
			"session_id", sessionID,
			"error", err,
		)
		return zero, usecase.ErrInternalError
	}

	return domain.NewSession(
		dto.ID,
		dto.UserID,
		dto.ExpiresAt,
	), nil
}

func (r *sessionRepo) Delete(ctx usecase.AppCtx, sessionID string) error {
	if err := r.client.Del(ctx.Context(), r.key(sessionID)).Err(); err != nil {
		r.log.Error(
			"не удалось удалить сессию из Redis",
			"session_id", sessionID,
			"error", err,
		)
		return usecase.ErrInternalError
	}

	return nil
}

func (r *sessionRepo) key(sessionID string) string {
	return "sessions:" + sessionID
}
