package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMain(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:8080/", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handleRequest(res, req)
	expected := ""
	actual := res.Body.String()
	if expected != actual {
		t.Fatalf("Expected: %s\nGot: %s.", expected, actual)
	}
}
