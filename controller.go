package goat

import (
	"net/http"
	"reflect"
)

type Controller interface {
	GetPath() string
	GetTags() []string
	GetDescription() string

	RequestSchemaSample() any
	ResponseSchemaSample() any

	MakeHandlerFunc(*Server) http.HandlerFunc
}

type Route[Props any, Return any] struct {
	Path             string // http.ServeMux path (eg. `GET /api/user/{id}`)
	Tags             []string
	Description      string
	PropDescriptions map[string]string // map of Prop names to descriptions
	Func             func(*Context[Props]) (status int, v *Return, err error)

	OverrideErrorHandler ErrorHandlerFunc
	OverrideEncoder      EncoderFunc

	blueprints []fieldBlueprint
}

func (r Route[Props, Return]) GetPath() string {
	return r.Path
}

func (r Route[Props, Return]) GetTags() []string {
	return r.Tags
}

func (r Route[Props, Return]) GetDescription() string {
	return r.Description
}

func (r Route[Props, Return]) RequestSchemaSample() any {
	return new(Props)
}

func (r Route[Props, Return]) ResponseSchemaSample() any {
	return new(Return)
}

func (route Route[Props, Return]) MakeHandlerFunc(s *Server) http.HandlerFunc {
	var sample Props
	blueprints, err := compileBlueprints(sample)
	if err != nil {
		s.logger.Error("got error while compiling blueprints (props) of", "controller", route.Path, "err", err)
		s.logger.Warn("this controller will not work as there has been an error", "controller", route.Path)
		return func(w http.ResponseWriter, r *http.Request) {}
	}
	route.blueprints = blueprints

	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("request", "controller", route.Path, "host", r.Host)
		// Reflect (if any blueprints exist)
		var reflection reflect.Value
		var props Props
		if len(route.blueprints) > 0 { // optimisation
			reflection = reflect.ValueOf(&props).Elem()
		}

		for _, b := range route.blueprints {
			err := b.SetField(reflection, s, r)
			if err != nil {
				s.logger.Error("error while parsing fields for", "controller", route.Path, "err", err)
				if route.OverrideErrorHandler != nil {
					route.OverrideErrorHandler(w, 400, err)
				} else {
					s.ErrorHandler(w, 400, err)
				}
				return
			}
		}

		// Run route
		ctx := Context[Props]{Props: props, Request: r, Writer: w, Logger: s.logger}
		status, v, err := route.Func(&ctx)

		if err != nil {
			s.logger.Error("controller returned an error", "controller", route.Path, "err", err)
			if route.OverrideErrorHandler != nil {
				route.OverrideErrorHandler(w, 400, err)
			} else {
				s.ErrorHandler(w, status, err)
			}
			return
		}

		// Respond
		w.WriteHeader(status)
		if v != nil {
			if route.OverrideEncoder != nil {
				route.OverrideEncoder(w, *v)
			} else {
				s.Encoder(w, *v)
			}
		}
	}
}
