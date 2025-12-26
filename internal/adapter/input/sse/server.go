package transportsse

import (
	"event-driven-architecture/internal/adapter/input/transport/sse/middleware/auth"
	"event-driven-architecture/internal/app"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	r   *gin.Engine
	app *app.App
}

func NewServer(app *app.App) *Server {
	r := gin.New()

	// --- middleware initialization ---
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))
	r.Use(auth.New(
		app.AppCtxManager,
		app.AuthenticateSessionUseCase,
	))

	// --- handler initialization ---
	handler := NewSSEHandler(app.Log, app.SSEBroker)

	// --- routes registration ---
	r.GET("/events", handler.Handle)

	return &Server{
		app: app,
	}
}

func (s *Server) Run(port int) error {
	s.app.Log.Info(fmt.Sprintf("Запуск SSE-сервера на порту %d", port))
	return s.r.Run(fmt.Sprintf(":%d", port))
}
