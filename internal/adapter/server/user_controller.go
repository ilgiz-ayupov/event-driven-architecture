package server

import (
	"event-driven-architecture/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	log               usecase.Logger
	createUserUsecase *usecase.CreateUserUseCase
}

func NewUserController(
	log usecase.Logger,
	createUserUsecase *usecase.CreateUserUseCase,
) *UserController {
	return &UserController{
		log:               log,
		createUserUsecase: createUserUsecase,
	}
}

func (e *UserController) RegisterRoutes(engine *gin.Engine) {
	g := engine.Group("/users")
	g.POST("", e.createUserHandler)
}

func (e *UserController) createUserHandler(c *gin.Context) {
	type Param struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var p Param
	if err := c.ShouldBindJSON(&p); err != nil {
		SendError(c, err)
		return
	}

	if err := e.createUserUsecase.Execute(c.Request.Context(), p.Email, p.Password); err != nil {
		SendError(c, err)
		return
	}

	SendData(c, http.StatusCreated, nil)
}
