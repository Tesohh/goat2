package goat

import (
	"net/http"
	"reflect"
)

type Controller interface {
	GetPath() string
	MakeHandlerFunc(*Server) http.HandlerFunc
}

type Route[Props any, Return any] struct {
	Path             string // http.ServeMux path (eg. `GET /api/user/{id}`)
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

func (route Route[Props, Return]) MakeHandlerFunc(s *Server) http.HandlerFunc {
	var sample Props
	route.blueprints = compileBlueprints(sample)

	return func(w http.ResponseWriter, r *http.Request) {
		// Reflect (if any blueprints exist)
		var reflection reflect.Value
		var props Props
		if len(route.blueprints) > 0 { // optimisation
			reflection = reflect.ValueOf(&props).Elem()
		}

		for _, b := range route.blueprints {
			err := b.SetField(reflection, s, r)
			if err != nil {
				if route.OverrideErrorHandler != nil {
					route.OverrideErrorHandler(w, 400, err)
				} else {
					s.ErrorHandler(w, 400, err)
				}
				return
			}
		}

		// Run route
		ctx := Context[Props]{Props: props, Request: r, Writer: w}
		status, v, err := route.Func(&ctx)

		if err != nil {
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
