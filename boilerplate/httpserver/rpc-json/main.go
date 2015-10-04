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

func main() {

	server := &jsonRpcServer{rpc.NewServer()}
	err := server.Register(&Service{})
	if err != nil {
		log.Fatal(err)
	}

	server.HandleHTTP("/rpc", "/debug")
	err = http.ListenAndServe(":5080", server)
	if err != nil {
		log.Println(err)
	}
}
