package goat

import (
	"net/http"

	"github.com/swaggest/openapi-go/openapi31"
)

type Server struct {
	mux         *http.ServeMux
	controllers []Controller

	reflector *openapi31.Reflector

	ErrorHandler ErrorHandlerFunc
	Encoder      EncoderFunc
}

// Returns a new pointer to a server, with the goat.DefaultErrorHandler and a goat.JSONEncoder.
func NewServer(info openapi31.Info) *Server {
	reflector := openapi31.NewReflector()
	reflector.Spec.Info = info

	return &Server{
		mux:          http.NewServeMux(),
		controllers:  make([]Controller, 0),
		reflector:    reflector,
		ErrorHandler: DefaultErrorHandler,
		Encoder:      JSONEncoder,
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
