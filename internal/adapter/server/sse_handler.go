package server

import (
	"event-driven-architecture/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (e *HTTPServer) sseHandler(c *gin.Context) {
	sessionID := c.MustGet("session_id").(string)

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

	subscription := e.broker.Subscribe(sessionID)
	defer subscription.Close()

	for {
		select {
		case <-r.Context().Done():
			return

		case msg, ok := <-subscription.Channel():
			e.log.Info(
				"SSE-хендлер: получение сообщения для отправки",
				"session_id", sessionID,
				"message_size", len(msg),
				"subscription_id", subscription.ID(),
				"active", ok,
				"msg", string(msg),
			)

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
