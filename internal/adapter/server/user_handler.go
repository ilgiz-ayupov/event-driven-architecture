package server

import (
	"event-driven-architecture/internal/domain"
	"event-driven-architecture/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (e *HTTPServer) createUserHandler(c *gin.Context) {
	type Param struct {
		Email    string `binding:"required"`
		Password string `binding:"required"`
	}

	var p Param
	if err := c.ShouldBindJSON(&p); err != nil {
		SendError(c, err)
		return
	}

	authUser := c.MustGet("auth_user").(domain.AuthUser)

	if err := Exec(
		c,
		e.transactionManager,
		e.appCtxManager,
		func(ctx usecase.AppCtx) error {
			return e.createUserUseCase.Execute(ctx, usecase.NewCreateUserInput(
				p.Email,
				p.Password,
				usecase.NewUserLog(authUser.ID),
			))
		},
	); err != nil {
		SendError(c, err)
		return
	}

	SendStatus(c, http.StatusCreated)
}
