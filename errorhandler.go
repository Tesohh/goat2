package goat

import (
	"fmt"
	"net/http"
)

type ErrorHandlerFunc func(http.ResponseWriter, int, error)

func DefaultErrorHandler(w http.ResponseWriter, status int, err error) {
	// TODO: Use JSON here.
	w.WriteHeader(status)
	fmt.Fprintf(w, "error: %s", err.Error())
}
