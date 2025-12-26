package transportsse

import (
	"event-driven-architecture/internal/usecase"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type SSEHandler struct {
	log    usecase.Logger
	broker usecase.EventBroker
}

func NewSSEHandler(log usecase.Logger, broker usecase.EventBroker) *SSEHandler {
	return &SSEHandler{
		log:    log,
		broker: broker,
	}
}

func (e *SSEHandler) Handle(c *gin.Context) {
	rawSessionID, ok := c.Get("session_id")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	sessionID, ok := rawSessionID.(string)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	e.log.Info(
		"SSE: подключение",
		"session_id", sessionID,
		"client_ip", c.ClientIP(),
	)
	defer e.log.Info(
		"SSE: отключение",
		"session_id", sessionID,
		"client_ip", c.ClientIP(),
	)

	w := c.Writer
	r := c.Request

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	flusher, ok := w.(http.Flusher)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	subscription := e.broker.Subscribe(sessionID)
	defer subscription.Close()

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			// клиент закрыл соединение
			return

		case <-ticker.C:
			// отправка комментария для поддержания соединения
			w.Write([]byte(": keep-alive\n\n"))
			flusher.Flush()

		case msg, ok := <-subscription.Channel():
			if !ok {
				// канал закрыт брокером
				return
			}

			// отправка SSE-сообщения
			if _, err := w.Write(msg); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}
