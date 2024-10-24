package main

// 1. Middleware
// 2. Routes defined per handler
// 3. InputFormDTO for POST/PATCH/DELETE requests automatically parsed and validated based on the handler
// 4. Dispatch valid input form DTO to handler
// 5. Response Object automatically parsed and validated and returned as JSON

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type (
	Action func(*Ctx) ActionResult
	Route  struct {
		Action   Action
		Method   string
		Path     string
		Request  Request
		Response Response
	}
	Request  interface{}
	Response interface{}

	Ctx struct {
		Request        *http.Request
		ResponseWriter http.ResponseWriter
		URLParams      httprouter.Params
	}
	ActionResult interface {
		Execute(*Ctx)
	}
)

type PlainActionResult struct {
	Body []byte
}

func (r *PlainActionResult) Execute(c *Ctx) {
	c.ResponseWriter.Header().Add("content-type", "text/plain")
	c.ResponseWriter.WriteHeader(200)
	c.ResponseWriter.Write(r.Body)
}

func home(c *Ctx) ActionResult {
	return &PlainActionResult{Body: []byte("Gharshana")}
}

func AddRoute(router *httprouter.Router, route Route) {
	router.Handle(route.Method, route.Path, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx := Ctx{
			ResponseWriter: w,
			Request:        r,
			URLParams:      params,
		}
		result := route.Action(&ctx)
		result.Execute(&ctx)
	})
}

func main() {
	log.Println("starting webfw")
	r := httprouter.New()
	AddRoute(r, Route{
		home,
		"GET",
		"/",
		nil,
		nil,
	})
	log.Fatalln(http.ListenAndServe(":8080", r))
}
