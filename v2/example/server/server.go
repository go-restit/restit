package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

// Server returns handler route the base URL, noun, nounp
type Server func(base, noun, nounp string) http.Handler

// New returns a new example server of a given type and store
func New(store *Store, f Factory) Server {

	Create := func(w http.ResponseWriter, r *http.Request) {
		if r == nil {
			http.Error(w, "request is nil", http.StatusBadRequest)
		}

		dec := json.NewDecoder(r.Body)
		item := f.Make().(Storable)
		if err := dec.Decode(item); err != nil {
			http.Error(w,
				fmt.Sprintf("failed to JSON decode request body as %T", item),
				http.StatusBadRequest)
		}

		store.Put(item)
	}

	Retrieve := func(w http.ResponseWriter, r *http.Request) {
	}

	return func(base, noun, nounp string) http.Handler {
		// path to use
		pathSingular := path.Join(base, noun, "{id}")
		pathPlural := path.Join(base, nounp)

		r := mux.NewRouter()
		r.HandleFunc(pathPlural, Create).Methods("POST")
		r.HandleFunc(pathSingular, Retrieve).Methods("GET")

		return r
	}

}
