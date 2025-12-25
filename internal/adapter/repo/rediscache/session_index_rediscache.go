package rediscache

import (
	"event-driven-architecture/internal/usecase"

	"github.com/redis/go-redis/v9"
)

type sessionIndexRepo struct {
	client *redis.Client
}

func NewSessionIndex(client *redis.Client) usecase.SessionIndexRepo {
	return &sessionIndexRepo{client: client}
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
		return err
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
		return nil
	default:
		return err
	}

	userKey := r.userSessionsKey(userID)

	pipe := r.client.TxPipeline()
	pipe.SRem(c, userKey, sessionID)
	pipe.Del(c, sessionKey)

	_, err = pipe.Exec(c)
	if err != nil {
		return err
	}

	return nil
}

func (r *sessionIndexRepo) SessionsByUser(ctx usecase.AppCtx, userID string) ([]string, error) {
	c := ctx.Context()

	userKey := r.userSessionsKey(userID)

	sessions, err := r.client.SMembers(c, userKey).Result()
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (r *sessionIndexRepo) userSessionsKey(userID string) string {
	return "user_sessions:" + userID
}

func (r *sessionIndexRepo) sessionUserKey(sessionID string) string {
	return "session_user:" + sessionID
}
