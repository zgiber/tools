package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/context"
)

type requestContext struct {
	context.Context
}

// convenience
func (ctx requestContext) request() *http.Request {
	if r, ok := ctx.Value("request").(*http.Request); ok {
		return r
	}

	return nil
}

func newContext(r *http.Request) requestContext {
	return requestContext{
		context.WithValue(context.Background(), "request", r),
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(r)
	fmt.Fprintln(w, ctx.request().UserAgent())
}

func main() {

	http.HandleFunc("/", handleRequest)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println(err)
	}

}
