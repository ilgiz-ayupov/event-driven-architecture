package server

import (
	"event-driven-architecture/internal/usecase"

	"github.com/gin-gonic/gin"
)

func AuthenticateSessionMiddleware(
	appCtxManager usecase.AppCtxManager,
	authenticateSessionUseCase *usecase.AuthenticateSessionUseCase,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("session_id")
		if err != nil {
			SendError(c, usecase.ErrUnauthorized)
			return
		}

		sessionID := cookie.Value

		appCtx, cancel := appCtxManager.CreateContext(c.Request.Context())
		defer cancel()

		authUser, err := authenticateSessionUseCase.Execute(appCtx, sessionID)
		if err != nil {
			SendError(c, err)
			return
		}

		c.Set("auth_user", authUser)
		c.Set("session_id", sessionID)
		c.Next()
	}
}
