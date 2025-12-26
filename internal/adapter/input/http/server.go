package transporthttp

import (
	v1 "event-driven-architecture/internal/adapter/input/transport/http/v1"
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
	engine := gin.New()

	// --- middleware initialization ---
	engine.Use(gin.Recovery())
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	// --- routes registration ---
	v1.RegisterRoutes(
		engine.Group("/v1"),
	)

	return &Server{
		r:   engine,
		app: app,
	}
}

func (s *Server) Run(port int) error {
	s.app.Log.Info("Запуск HTTP-сервера на порту", port)
	return s.r.Run(fmt.Sprintf(":%d", port))
}
