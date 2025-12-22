package appctx

import "errors"

var (
	errSessionNotFound = errors.New("сессия не найдена в контексте")
)
