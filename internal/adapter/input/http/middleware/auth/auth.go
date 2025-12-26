package auth

import (
	"event-driven-architecture/internal/domain"
	"event-driven-architecture/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	contextAuthUserKey  = "auth_user"
	contextSessionIDKey = "session_id"
)

func GetAuthUser(c *gin.Context) (domain.AuthUser, bool) {
	raw, exists := c.Get(contextAuthUserKey)
	if !exists {
		return domain.AuthUser{}, false
	}
	user, ok := raw.(domain.AuthUser)
	return user, ok
}

func GetSessionID(c *gin.Context) (string, bool) {
	raw, exists := c.Get(contextSessionIDKey)
	if !exists {
		return "", false
	}
	sessionID, ok := raw.(string)
	return sessionID, ok
}

func New(
	appCtxManager usecase.AppCtxManager,
	uc *usecase.AuthenticateSessionUseCase,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("session_id")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		sessionID := cookie.Value

		appCtx, cancel := appCtxManager.CreateContext(c.Request.Context())
		authUser, err := uc.Execute(appCtx, sessionID)
		if err != nil {
			cancel()
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		cancel()

		c.Set(contextAuthUserKey, authUser)
		c.Set(contextSessionIDKey, sessionID)

		c.Next()
	}
}
