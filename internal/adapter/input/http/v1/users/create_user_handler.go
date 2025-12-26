package users

import (
	transporthttp "event-driven-architecture/internal/adapter/input/transport/http"
	"event-driven-architecture/internal/adapter/input/transport/http/middleware/auth"
	"event-driven-architecture/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateUserHandler struct {
	txManager     usecase.TransactionManager
	appCtxManager usecase.AppCtxManager
	uc            *usecase.CreateUserUseCase
}

func NewCreateUserHandler(
	txManager usecase.TransactionManager,
	appCtxManager usecase.AppCtxManager,
	uc *usecase.CreateUserUseCase,
) *CreateUserHandler {
	return &CreateUserHandler{
		txManager:     txManager,
		appCtxManager: appCtxManager,
		uc:            uc,
	}
}

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type CreateUserResponse struct {
	UserID string `json:"user_id"`
}

func (h *CreateUserHandler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transporthttp.SendError(c, err)
		return
	}

	authUser, ok := auth.GetAuthUser(c)
	if !ok {
		transporthttp.SendError(c, usecase.ErrUnauthorized)
		return
	}

	out, err := transporthttp.ReturnWithCommit(
		c.Request.Context(),
		h.txManager,
		h.appCtxManager,
		func(ctx usecase.AppCtx) (usecase.CreateUserOutput, error) {
			return h.uc.Execute(ctx, usecase.NewCreateUserInput(
				req.Email,
				req.Password,
				usecase.NewUserLog(authUser.ID),
			))
		},
	)
	if err != nil {
		transporthttp.SendError(c, err)
		return
	}

	transporthttp.SendData(c, http.StatusCreated, CreateUserResponse{UserID: out.UserID})
}
