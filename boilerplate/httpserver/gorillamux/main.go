package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-origin", "*")
	fmt.Fprintf(w, "Welcome %s, %s\n", r.RemoteAddr, r.UserAgent())
	return
}

func GETProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-origin", "*")
	enc := json.NewEncoder(w)
	err := enc.Encode(map[string]interface{}{
		"id":    1,
		"name":  "Bicycle",
		"color": "red",
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GETArticlesCategoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-origin", "*")
	enc := json.NewEncoder(w)
	err := enc.Encode(map[string]interface{}{
		"requestURL": r.URL.String(),
		"vars":       mux.Vars(r),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GETArticleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-origin", "*")
	enc := json.NewEncoder(w)
	err := enc.Encode(map[string]interface{}{
		"requestURL": r.URL.String(),
		"vars":       mux.Vars(r),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func POSTProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-origin", "*")

	type product struct {
		ID       int
		Name     string
		Category string
	}

	prod := &product{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(prod)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	enc := json.NewEncoder(w)
	err = enc.Encode(prod)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	r := mux.NewRouter()
	get := r.Methods("GET").Subrouter()
	post := r.Methods("POST").Subrouter()

	get.HandleFunc("/", HomeHandler)                                      // Simple handler
	get.HandleFunc("/products", GETProductHandler)                        //
	post.HandleFunc("/products", POSTProductHandler)                      // POST example
	get.HandleFunc("/articles/{category}", GETArticlesCategoryHandler)    // Variable
	get.HandleFunc("/articles/{category}/{id:[0-9]+}", GETArticleHandler) // Variable with regexp match

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Println(err)
	}
}
