package main

import (
	"github.com/Tesohh/goat"
	"github.com/swaggest/openapi-go/openapi31"
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

type HelloReq struct {
	Name string `goat:",path" path:"name"`
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

	s := goat.NewServer(openapi31.Info{
		Title: "Greeter API",
	})
	s.AddController(route)

	err := s.CompileOpenAPI()
	if err != nil {
		panic(err)
	}
	s.AddSwaggerUI()

	s.Listen(":8080")
}
