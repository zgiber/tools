package main

import (
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"

	"golang.org/x/net/websocket"
)

func serveRPC(ws *websocket.Conn) {
	jsonrpc.ServeConn(ws)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	fmt.Fprint(w, "<!DOCTYPE html>Welcome!")
}

func main() {

	err := rpc.Register(&Service{})
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", mainHandler) // empty page for testing via browser (required by origin check in websockets)
	http.Handle("/ws", websocket.Handler(serveRPC))

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
