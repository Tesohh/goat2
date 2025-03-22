# goat
self documenting backend framework for go

---

`goat` is a backend framework for go, based on `net/http`, 
that automatically generates OpenAPI 3.1 specifications and is a joy to use.

It uses reflection under the hood, 
although "request-time" reflection is kept minimal; 
most work is done at "boot-time".

## why goat?
I was mesmerized by ASP.NET Core's ability to automatically generate Swagger specs, 
manage query/body/path params automatically, 
but I didn't want to cheat on Go by using C#, 
which I personally don't like writing, so here we are.

## Example
```go
type HelloReq struct {
	Name string `path:"name"`
	Surname string `query:"surname"`
	HelloBody
}

type HelloBody struct {
	ID string `json:"id"`
}

type HelloRes struct {
	Name string `json:"name"`
}

// for more information on struct tags, 
// refer to (swaggest/openapi-go)[https://github.com/swaggest/openapi-go]
// note that for now, the only "sources" are `path`, `query` and body.

func main() {
	s := goat.NewServer(openapi31.Info{
		Title: "Greeter API",
	})

	route := goat.Route[HelloReq, HelloRes]{
		Path: "POST /hello/{name}", // http.ServeMux-like path
		Func: func(c *goat.Context[HelloReq]) (int, *HelloRes, error) {
			return 100, &HelloRes{
				Name: c.Props.Name,
			}, nil
		},
	}

	s.AddController(route)

	err := s.CompileOpenAPI()
	if err != nil {
		panic(err)
	}

	s.AddSwaggerUI(swgui.Config{})

	s.Listen(":8080")
}

```

## features to come
- [ ] Middleware
