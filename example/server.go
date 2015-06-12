package example

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"text/template"
)

var txtList string = `{
	"code": 302,
	"status": "success",
	"nodes": [
		{
			"id": 1,
			"name": "node 1",
			"desc": "example node 1"
		},
		{
			"id": 2,
			"name": "node 2",
			"desc": "example node 2"
		},
		{
			"id": 3,
			"name": "node 3",
			"desc": "example node 3"
		}
	]
}`

var txtRead string = `{
	"code": 302,
	"status": "success",
	"nodes": [
		{
			"id": {{ .Vars.id }},
			"name": "node {{ .Vars.id }}",
			"desc": "example node {{ .Vars.id }}"
		}
	]
}`

var txtWrite string = `{
	"code": 302,
	"status": "success",
	"nodes": [
		{
			"id": 4,
			"name": "node 4",
			"desc": "example node 4"
		}
	]
}`

var txtUpdate string = `{
	"code": 302,
	"status": "success",
	"nodes": [
		{
			"id": {{ .Vars.id }},
			"name": "node {{ .Vars.id }} updated",
			"desc": "example node {{ .Vars.id }} updated"
		}
	]
}`

var txtDelete string = `{
	"code": 404,
	"status": "success",
	"nodes": [
		{
			"id": {{ .Vars.id }},
			"name": "node {{ .Vars.id }} updated",
			"desc": "example node {{ .Vars.id }} updated"
		}
	]
}`

var txtStatusMethodNotAllowed string = `{
	"code": 405,
	"status": "error",
	"message": "Method not allowed"
}`

// ExampleHandler is a dummy ordinarly REST server.
// It doesn't create or delete anything. Only returns
// static responses.
func ExampleHandler() http.Handler {
	r := mux.NewRouter()

	tplRead := template.Must(template.New("read").Parse(txtRead))
	tplUpdate := template.Must(template.New("update").Parse(txtUpdate))
	tplDelete := template.Must(template.New("delete").Parse(txtDelete))

	// handles listing and create
	r.HandleFunc("/api/nodes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			// dummy listing response
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, txtList)
			return
		} else if r.Method == "POST" {
			// dummy node creating response
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, txtWrite)
			return
		}
		// no other method allowed
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, txtStatusMethodNotAllowed)
	})

	// handles read, update, delete
	r.HandleFunc("/api/nodes/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			tplRead.Execute(w, map[string]interface{}{
				"Vars": vars,
			})
		} else if r.Method == "PUT" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)
			tplUpdate.Execute(w, map[string]interface{}{
				"Vars": vars,
			})
			return
		} else if r.Method == "DELETE" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			tplDelete.Execute(w, map[string]interface{}{
				"Vars": vars,
			})
			return
		}
	})

	return r
}
