package multi

import (
	"event-driven-architecture/internal/domain"
	"event-driven-architecture/internal/usecase"
)

type multiSessionRepo struct {
	sessionRepos []usecase.SessionRepo
}

func NewMultiSession(repos ...usecase.SessionRepo) usecase.SessionRepo {
	return &multiSessionRepo{
		sessionRepos: repos,
	}
}

func (r *multiSessionRepo) Create(ctx usecase.AppCtx, session domain.Session) error {
	for _, repo := range r.sessionRepos {
		if err := repo.Create(ctx, session); err != nil {
			return err
		}
	}

	return nil
}

func (r *multiSessionRepo) Find(ctx usecase.AppCtx, sessionID string) (domain.Session, error) {
	var session domain.Session
	var err error

	for _, repo := range r.sessionRepos {
		session, err = repo.Find(ctx, sessionID)
		if err == nil {
			return session, nil
		}
	}

	return domain.Session{}, err
}

func (r *multiSessionRepo) Delete(ctx usecase.AppCtx, sessionID string) error {
	for _, repo := range r.sessionRepos {
		if err := repo.Delete(ctx, sessionID); err != nil {
			return err
		}
	}

	return nil
}
