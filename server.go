package goat

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/swaggest/openapi-go/openapi31"
	"github.com/swaggest/swgui"
	"github.com/swaggest/swgui/v5cdn"
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

	fullPath := c.GetPath()
	s.mux.Handle(fullPath, c.MakeHandlerFunc(s))

	parts := strings.Fields(fullPath)
	method := ""
	path := ""
	if len(parts) == 1 {
		method = "GET"
		path = parts[0]
	} else if len(parts) == 2 {
		method = parts[0]
		path = parts[1]
	} else {
		panic(fmt.Sprintf("unable to correctly split parts of %s. Can have 1 or 2 parts (method and path)", fullPath)) // TODO: switch to slogging
	}

	ctx, err := s.reflector.NewOperationContext(method, path)
	if err != nil {
		panic(err)
	}

	ctx.SetTags(c.GetTags()...)
	ctx.SetDescription(c.GetDescription())
	ctx.AddReqStructure(c.RequestSchemaSample())
	ctx.AddRespStructure(c.ResponseSchemaSample())

	err = s.reflector.AddOperation(ctx)
	if err != nil {
		panic(err)
	}
}

func (s *Server) CompileOpenAPI() error {
	b, err := s.reflector.Spec.MarshalJSON()
	if err != nil {
		return err
	}

	s.mux.HandleFunc("/api/docs/v3.1/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write(b)
	})

	return nil
}

// NOTE: you must have run Server.CompileOpenAPI beforehand.
func (s *Server) AddSwaggerUI(config swgui.Config) error {
	ui := v5cdn.NewWithConfig(config)(s.reflector.Spec.Info.Title, "/api/docs/v3.1/openapi.json", "/api/docs")
	s.mux.Handle("/api/docs", ui)

	return nil
}

func (s *Server) Listen(addr string) {
	http.ListenAndServe(addr, s.mux)
}
