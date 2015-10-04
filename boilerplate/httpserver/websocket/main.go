package main

import (
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

// Echo the data received on the WebSocket.
func EchoServer(ws *websocket.Conn) {
	io.Copy(ws, ws)
}

// This example demonstrates a trivial echo server.
func main() {

	err := http.ListenAndServe(":8080", websocket.Handler(EchoServer))
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
