package main

import (
	"io"
	"log"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

var connected = "200 Connected to Go RPC"

// serve rpc-json request over http
// alternatively we can use rpc.
type jsonRpcServer struct {
	*rpc.Server
}

func (handler *jsonRpcServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "CONNECT" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "405 must CONNECT\n")
		return
	}
	conn, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		log.Print("rpc hijacking ", req.RemoteAddr, ": ", err.Error())
		return
	}
	io.WriteString(conn, "HTTP/1.0 "+connected+"\n\n")
	handler.ServeCodec(jsonrpc.NewServerCodec(conn))
}

type Api struct{}

func (a *Api) DoSomething(args interface{}, resp *int) error {
	log.Println(args)
	*resp = 42
	return nil
}

func main() {

	server := &jsonRpcServer{rpc.NewServer()}
	err := server.Register(&Api{}) //in service.go
	if err != nil {
		log.Fatal(err)
	}

	server.HandleHTTP("/", "/debug")
	err = http.ListenAndServe(":5080", server)
	if err != nil {
		log.Println(err)
	}
}
