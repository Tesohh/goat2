package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Tesohh/goat"
	"github.com/lmittmann/tint"
	"github.com/swaggest/openapi-go/openapi31"
	"github.com/swaggest/swgui"
)

type Person struct {
	Name string
}

//
// type World struct {
// 	People []Person
// }
//
// type Hello struct {
// 	World
// }

// type Hello struct {
// 	Name   string `goat:",path" path:"name"`
// 	Person `goat:"person"`
// }

type HelloBody struct {
	ID string `json:"id"`
}

type HelloReq struct {
	Name string `path:"name"`
	HelloBody
}

type HelloRes struct {
	Name string `json:"name"`
}

func main() {
	route := goat.Route[HelloReq, HelloRes]{
		Path:        "/hello/{name}",
		Description: "hello",
		Func: func(c *goat.Context[HelloReq]) (int, *HelloRes, error) {
			return 100, &HelloRes{
				Name: c.Props.Name,
			}, nil
		},
	}
	routePost := goat.Route[HelloReq, HelloRes]{
		Path: "POST /hello/{name}",
		Func: func(c *goat.Context[HelloReq]) (int, *HelloRes, error) {
			return 100, &HelloRes{
				Name: c.Props.Name,
			}, nil
		},
	}
	errorRoute := goat.Route[struct{}, struct{}]{
		Path: "GET /error",
		Func: func(c *goat.Context[struct{}]) (int, *struct{}, error) {
			c.Logger.Info("hi from error route")
			return 100, &struct{}{}, fmt.Errorf("errorrr")
		},
	}

	s := goat.NewServer(openapi31.Info{
		Title: "Greeter API",
	})
	s.SetLogger(slog.New(tint.NewHandler(os.Stderr, nil)))

	s.AddController(route)
	s.AddController(routePost)
	s.AddController(errorRoute)

	err := s.CompileOpenAPI()
	if err != nil {
		panic(err)
	}

	config := swgui.Config{
		SettingsUI: map[string]string{
			"tryItOutEnabled": "true",
		},
	}

	s.AddSwaggerUI(config)

	s.Listen(":8080")
}
