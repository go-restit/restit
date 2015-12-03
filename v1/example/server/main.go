package main

import (
	"github.com/yookoala/restit/example"
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
