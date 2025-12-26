package transporthttp

import (
	"event-driven-architecture/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuccessResponse представляет успешный ответ API
type SuccessResponse struct {
	Data any `json:"data"`
}

// ErrorResponse представляет ответ с ошибкой API
type ErrorResponse struct {
	Error string `json:"error"`
}

// SendData отправляет статус код
func SendStatus(c *gin.Context, statusCode int) {
	c.Status(statusCode)
}

// SendData отправляет успешный ответ
func SendData(c *gin.Context, code int, data any) {
	c.JSON(code, SuccessResponse{Data: data})
}

// SendError отправляет ответ с ошибкой
func SendError(c *gin.Context, err error) {
	code := mapErrorToStatus(err)

	c.JSON(code, ErrorResponse{Error: err.Error()})
	c.Abort()
}

func mapErrorToStatus(err error) (code int) {
	switch err {
	case usecase.ErrNoData:
		return http.StatusNotFound

	case usecase.ErrUnauthorized:
		return http.StatusUnauthorized

	case usecase.ErrInternalError:
		return http.StatusInternalServerError

	default:
		return http.StatusBadRequest
	}
}
