package rediscache

import (
	"event-driven-architecture/internal/usecase"

	"github.com/redis/go-redis/v9"
)

type sessionIndexRepo struct {
	log    usecase.Logger
	client *redis.Client
}

func NewSessionIndex(log usecase.Logger, client *redis.Client) usecase.SessionIndexRepo {
	return &sessionIndexRepo{log: log, client: client}
}

func (r *sessionIndexRepo) Add(ctx usecase.AppCtx, sessionID, userID string) error {
	c := ctx.Context()

	userKey := r.userSessionsKey(userID)
	sessionKey := r.sessionUserKey(sessionID)

	pipe := r.client.TxPipeline()
	pipe.SAdd(c, userKey, sessionID)
	pipe.Set(c, sessionKey, userID, 0)

	_, err := pipe.Exec(c)
	if err != nil {
		r.log.Error(
			"не удалось индексировать сессию пользователя в Redis",
			"session_id", sessionID,
			"user_id", userID,
			"error", err,
		)
		return usecase.ErrInternalError
	}

	return nil
}

func (r *sessionIndexRepo) Remove(ctx usecase.AppCtx, sessionID string) error {
	c := ctx.Context()

	sessionKey := r.sessionUserKey(sessionID)

	userID, err := r.client.Get(c, sessionKey).Result()
	switch err {
	case nil:
	case redis.Nil:
		r.log.Warn(
			"сессия не найдена в индексе Redis",
			"session_id", sessionID,
		)
		return nil
	default:
		r.log.Error(
			"не удалось получить сессию из индекса Redis",
			"session_id", sessionID,
			"error", err,
		)
		return usecase.ErrInternalError
	}

	userKey := r.userSessionsKey(userID)

	pipe := r.client.TxPipeline()
	pipe.SRem(c, userKey, sessionID)
	pipe.Del(c, sessionKey)

	_, err = pipe.Exec(c)
	if err != nil {
		r.log.Error(
			"не удалось удалить сессию из индекса Redis",
			"session_id", sessionID,
			"user_id", userID,
			"error", err,
		)
		return usecase.ErrInternalError
	}

	return nil
}

func (r *sessionIndexRepo) SessionsByUser(ctx usecase.AppCtx, userID string) ([]string, error) {
	c := ctx.Context()

	userKey := r.userSessionsKey(userID)

	sessions, err := r.client.SMembers(c, userKey).Result()
	if err != nil {
		r.log.Error(
			"не удалось получить сессии пользователя из индекса Redis",
			"user_id", userID,
			"error", err,
		)
		return nil, usecase.ErrInternalError
	}

	return sessions, nil
}

func (r *sessionIndexRepo) userSessionsKey(userID string) string {
	return "user_sessions:" + userID
}

func (r *sessionIndexRepo) sessionUserKey(sessionID string) string {
	return "session_user:" + sessionID
}
