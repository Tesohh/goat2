package main

import (
	"github.com/Tesohh/goat"
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

	s := goat.NewServer(openapi31.Info{
		Title: "Greeter API",
	})
	s.AddController(route)
	s.AddController(routePost)

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
