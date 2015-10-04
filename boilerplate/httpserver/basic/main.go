package main

import (
	"log"
	"net/http"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {

}

func main() {

	http.HandleFunc("/", handleRequest)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println(err)
	}

}
