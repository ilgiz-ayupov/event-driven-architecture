package server

import (
	"context"
	"event-driven-architecture/internal/usecase"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	log                usecase.Logger
	address            string
	transactionManager usecase.TransactionManager
	appCtxManager      usecase.AppCtxManager
	broker             usecase.EventBroker

	authenticateSessionUseCase *usecase.AuthenticateSessionUseCase
	loginUserUseCase           *usecase.LoginUserUseCase
	createUserUseCase          *usecase.CreateUserUseCase
}

func NewHTTPServer(
	log usecase.Logger,
	port int,
	transactionManager usecase.TransactionManager,
	appCtxManager usecase.AppCtxManager,
	broker usecase.EventBroker,

	authenticateSessionUseCase *usecase.AuthenticateSessionUseCase,
	loginUserUseCase *usecase.LoginUserUseCase,
	createUserUseCase *usecase.CreateUserUseCase,
) *HTTPServer {
	return &HTTPServer{
		log:                log,
		address:            fmt.Sprintf(":%d", port),
		transactionManager: transactionManager,
		appCtxManager:      appCtxManager,
		broker:             broker,

		authenticateSessionUseCase: authenticateSessionUseCase,
		loginUserUseCase:           loginUserUseCase,
		createUserUseCase:          createUserUseCase,
	}
}

func (e *HTTPServer) Run(ctx context.Context) {
	engine := gin.New()

	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	// инициализировать middleware
	authenticateSessionMiddleware := AuthenticateSessionMiddleware(
		e.appCtxManager,
		e.authenticateSessionUseCase,
	)

	// регистрация путей
	apiV1 := engine.Group("/api/v1")

	authGroup := apiV1.Group("/auth")
	authGroup.POST("/login", e.login)

	userGroup := apiV1.Group("/users", authenticateSessionMiddleware)
	userGroup.POST("", e.createUserHandler)

	sseGroup := apiV1.Group("/events", authenticateSessionMiddleware)
	sseGroup.GET("", e.sseHandler)

	// создание сервера
	srv := &http.Server{
		Addr:    e.address,
		Handler: engine,
	}

	// обработка сигнала завершения работы
	go func() {
		<-ctx.Done()
		e.log.Info("остановка http-сервера...", "address", e.address)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			e.log.Error("ошибка при остановке http-сервера:", "error", err)
			return
		}
		e.log.Info("http-сервер остановлен", "address", e.address)
	}()

	// запуск сервера
	e.log.Info("http-сервер запущен", "address", e.address)
	srv.ListenAndServe()
}
