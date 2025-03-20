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

type Hello struct {
	Name string `goat:",path"`
	Person
}

func main() {
	route := goat.Route[Hello, Hello]{
		Path:        "/hello/{name}",
		Description: "hello",
		Func: func(*goat.Context[Hello]) (int, *Hello, error) {
			return 100, &Hello{
				Name:   "AIOEIOFJIOJWEOFJ",
				Person: Person{Name: "TESTING"},
			}, nil
		},
	}

	s := goat.NewServer(openapi31.Info{})
	s.AddController(route)
	s.Listen(":8080")
}
