package server

import (
	"event-driven-architecture/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuccessResponse представляет успешный ответ API
type SuccessResponse struct {
	Data any `json:",omitempty"`
}

// ErrorResponse представляет ответ с ошибкой API
type ErrorResponse struct {
	Error string
}

// SendData отправляет статус код
func SendStatus(c *gin.Context, statusCode int) {
	c.Status(statusCode)
}

// SendData отправляет успешный ответ
func SendData(c *gin.Context, statusCode int, data *any) {
	response := SuccessResponse{Data: data}
	c.JSON(statusCode, response)
}

// SendError отправляет ответ с ошибкой
func SendError(c *gin.Context, err error) {
	statusCode := mapErrorToStatus(err)
	response := ErrorResponse{Error: err.Error()}

	c.JSON(statusCode, response)
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
