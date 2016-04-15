package main

import (
	"github.com/go-restit/restit/v1/example"
	"net/http"
	"os"
)

// port to use
var port string

func init() {
	port = "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
}

func main() {
	http.ListenAndServe(":"+port, example.ExampleHandler())
}
