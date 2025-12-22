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
	log     usecase.Logger
	address string
	broker  usecase.EventBroker

	createUserUseCase *usecase.CreateUserUseCase
}

func NewHTTPServer(
	log usecase.Logger,
	port int,
	broker usecase.EventBroker,

	createUserUseCase *usecase.CreateUserUseCase,
) *HTTPServer {
	return &HTTPServer{
		log:     log,
		broker:  broker,
		address: fmt.Sprintf(":%d", port),

		createUserUseCase: createUserUseCase,
	}
}

func (e *HTTPServer) Run(ctx context.Context) {
	engine := gin.New()

	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// регистрация путей
	e.userController().RegisterRoutes(engine)
	e.sseController().RegisterRoutes(engine)

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

func (e *HTTPServer) userController() *UserController {
	return NewUserController(
		e.log,
		e.createUserUseCase,
	)
}

func (e *HTTPServer) sseController() *SSEController {
	return NewSSEController(
		e.log,
		e.broker,
	)
}
