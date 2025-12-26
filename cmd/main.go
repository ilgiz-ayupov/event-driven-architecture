package main

import (
	transporthttp "event-driven-architecture/internal/adapter/input/transport/http"
	transportsse "event-driven-architecture/internal/adapter/input/transport/sse"
	"event-driven-architecture/internal/app"
)

func main() {
	cfg := app.LoadConfig()

	app, err := app.Build(cfg)
	if err != nil {
		panic(err)
	}

	// запуск SSE-сервера
	sseServer := transportsse.NewServer(app)
	go func() {
		if err := sseServer.Run(cfg.SSE.Port); err != nil {
			panic(err)
		}
	}()

	// запуск HTTP-сервера
	httpServer := transporthttp.NewServer(app)
	go func() {
		if err := httpServer.Run(cfg.HTTP.Port); err != nil {
			panic(err)
		}
	}()
}
