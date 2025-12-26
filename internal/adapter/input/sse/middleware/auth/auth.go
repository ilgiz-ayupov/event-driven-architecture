package auth

import (
	"event-driven-architecture/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

		c.Set("auth_user", authUser)
		c.Set("session_id", sessionID)

		c.Next()
	}
}
