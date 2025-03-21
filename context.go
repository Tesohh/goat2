package goat

import (
	"log/slog"
	"net/http"
)

type Context[T any] struct {
	Props   T
	Request *http.Request
	Writer  http.ResponseWriter

	Logger *slog.Logger
}
