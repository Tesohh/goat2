package goat

import (
	"net/http"

	"github.com/swaggest/openapi-go/openapi31"
)

type Server struct {
	mux         *http.ServeMux
	controllers []Controller

	ErrorHandler ErrorHandlerFunc
	Encoder      EncoderFunc
	Info         openapi31.Info
}

// Returns a new pointer to a server, with the goat.DefaultErrorHandler and a goat.JSONEncoder.
func NewServer(info openapi31.Info) *Server {
	return &Server{
		mux:          http.NewServeMux(),
		controllers:  make([]Controller, 0),
		ErrorHandler: DefaultErrorHandler,
		Encoder:      JSONEncoder,
		Info:         info,
	}
}

func (s *Server) AddController(c Controller) {
	s.controllers = append(s.controllers, c)

	path := c.GetPath()
	s.mux.Handle(path, c.MakeHandlerFunc(s))
}

func (s *Server) Listen(addr string) {
	http.ListenAndServe(addr, s.mux)
}
