package server

import (
	"event-driven-architecture/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SSEController struct {
	log    usecase.Logger
	broker usecase.EventBroker
}

func NewSSEController(
	log usecase.Logger,
	broker usecase.EventBroker,
) *SSEController {
	return &SSEController{
		log:    log,
		broker: broker,
	}
}

func (e *SSEController) RegisterRoutes(engine *gin.Engine) {
	engine.GET("/events", e.handler)
}

func (e *SSEController) handler(c *gin.Context) {
	w := c.Writer
	r := c.Request

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		SendError(c, usecase.ErrInternalError)
		return
	}

	subscriberID, ch := e.broker.Subscribe()
	defer e.broker.Unsubscribe(subscriberID)

	for {
		select {
		case <-r.Context().Done():
			return

		case msg, ok := <-ch:
			if !ok {
				// канал закрыт брокером
				return
			}

			// отправка SSE-сообщения
			w.Write(msg)
			w.Write([]byte("\n\n"))
			flusher.Flush()
		}
	}
}
