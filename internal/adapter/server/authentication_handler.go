package server

import (
	"event-driven-architecture/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (e *HTTPServer) login(c *gin.Context) {
	type Param struct {
		Email    string `binding:"required"`
		Password string `binding:"required"`
	}

	var p Param
	if err := c.ShouldBindJSON(&p); err != nil {
		SendError(c, err)
		return
	}

	sessionID, err := ReturnWithCommit(
		c,
		e.transactionManager,
		e.appCtxManager,
		func(ctx usecase.AppCtx) (string, error) {
			return e.loginUserUseCase.Execute(ctx, p.Email, p.Password)
		})
	if err != nil {
		SendError(c, err)
		return
	}

	// set cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteNoneMode,
	})

	SendStatus(c, http.StatusOK)
}
