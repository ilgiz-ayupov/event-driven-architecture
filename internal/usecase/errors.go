package usecase

import "errors"

var (
	// ErrNoData нет данных
	ErrNoData = errors.New("нет данных")
	// ErrInvalidCredentials недействительные учетные данные
	ErrInvalidCredentials = errors.New("недействительные учетные данные")
	// ErrUnauthorized
	ErrUnauthorized = errors.New("неавторизованный")
	// ErrInternalError внутрення ошибка
	ErrInternalError = errors.New("внутренняя ошибка")
)
