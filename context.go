package goat

import "net/http"

type Context[T any] struct {
	Props   T
	Request *http.Request
	Writer  http.ResponseWriter
}
