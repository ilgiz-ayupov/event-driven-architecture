package auth

import (
	transporthttp "event-driven-architecture/internal/adapter/input/transport/http"
	"event-driven-architecture/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthLoginHandler struct {
	txManager     usecase.TransactionManager
	appCtxManager usecase.AppCtxManager
	uc            *usecase.LoginUserUseCase
}

func NewAuthLoginHandler(
	txManager usecase.TransactionManager,
	appCtxManager usecase.AppCtxManager,
	uc *usecase.LoginUserUseCase,
) *AuthLoginHandler {
	return &AuthLoginHandler{
		txManager:     txManager,
		appCtxManager: appCtxManager,
		uc:            uc,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	SessionID string `json:"session_id"`
}

func (h *AuthLoginHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transporthttp.SendError(c, err)
		return
	}

	sessionID, err := transporthttp.ReturnWithCommit(
		c.Request.Context(),
		h.txManager,
		h.appCtxManager,
		func(ctx usecase.AppCtx) (string, error) {
			return h.uc.Execute(ctx, req.Email, req.Password)
		})
	if err != nil {
		transporthttp.SendError(c, err)
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

	transporthttp.SendData(c, http.StatusOK, LoginResponse{SessionID: sessionID})
}
