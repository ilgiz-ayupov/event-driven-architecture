package v1

import (
	"event-driven-architecture/internal/adapter/input/transport/http/v1/auth"
	"event-driven-architecture/internal/adapter/input/transport/http/v1/users"
	"event-driven-architecture/internal/usecase"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	rg *gin.RouterGroup,
	txManager usecase.TransactionManager,
	appCtxManger usecase.AppCtxManager,
	loginUserUseCase *usecase.LoginUserUseCase,
	createUserUseCase *usecase.CreateUserUseCase,
) {
	loginHandler := auth.NewAuthLoginHandler(
		txManager,
		appCtxManger,
		loginUserUseCase,
	)

	createUserHandler := users.NewCreateUserHandler(
		txManager,
		appCtxManger,
		createUserUseCase,
	)

	auth := rg.Group("/auth")
	auth.POST("/login", loginHandler.Login)

	users := rg.Group("/users")
	users.POST("", createUserHandler.Create)
}
