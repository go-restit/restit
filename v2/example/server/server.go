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

	return func(base, noun, nounp string) http.Handler {

		List := func(w http.ResponseWriter, r *http.Request) {

			if r == nil {
				http.Error(w, "request is nil", http.StatusBadRequest)
			}

			list := store.List(noun)

			w.WriteHeader(http.StatusOK)
			enc := json.NewEncoder(w)
			enc.Encode(map[string]interface{}{
				"status": http.StatusOK,
				nounp:    list,
			})
		}

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

			w.WriteHeader(http.StatusOK)
			enc := json.NewEncoder(w)
			enc.Encode(map[string]interface{}{
				"status": http.StatusOK,
				nounp: []interface{}{
					item,
				},
				noun: item,
			})
		}

		Retrieve := func(w http.ResponseWriter, r *http.Request) {

			// search the old stored item
			// if not found, return 404
			id := mux.Vars(r)["id"]

			item := store.Get(noun, id)
			enc := json.NewEncoder(w)

			if item == nil {
				w.WriteHeader(http.StatusNotFound)
				enc.Encode(map[string]interface{}{
					"status":  http.StatusNotFound,
					"message": "not found",
				})
				return
			}

			w.WriteHeader(http.StatusOK)
			enc.Encode(map[string]interface{}{
				"status": http.StatusOK,
				nounp: []interface{}{
					item,
				},
				noun: item,
			})
		}

		Update := func(w http.ResponseWriter, r *http.Request) {

			if r == nil {
				http.Error(w, "request is nil", http.StatusBadRequest)
			}

			// decode the request body
			dec := json.NewDecoder(r.Body)
			item := f.Make().(Storable)
			if err := dec.Decode(item); err != nil {
				http.Error(w,
					fmt.Sprintf("failed to JSON decode request body as %T", item),
					http.StatusBadRequest)
			}

			// encoder for returning protocol
			enc := json.NewEncoder(w)

			// search the old stored item
			// if not found, return 404
			id := mux.Vars(r)["id"]
			if v := store.Get(item.GetType(), id); v == nil {
				w.WriteHeader(http.StatusNotFound)
				enc.Encode(map[string]interface{}{
					"status":  http.StatusNotFound,
					"message": "not found",
				})
				return
			}

			// enforece ID, then put it into the store
			item.SetID(id)
			store.Put(item)

			w.WriteHeader(http.StatusOK)
			enc.Encode(map[string]interface{}{
				"status": http.StatusOK,
				nounp: []interface{}{
					item,
				},
				noun: item,
			})
		}

		Patch := func(w http.ResponseWriter, r *http.Request) {
			type storablePatchable interface {
				Storable
				Patchable
			}

			if r == nil {
				http.Error(w, "request is nil", http.StatusBadRequest)
			}

			// Decode the request body.
			dec := json.NewDecoder(r.Body)
			item := f.Make().(Storable)
			if err := dec.Decode(item); err != nil {
				http.Error(w,
					fmt.Sprintf("failed to JSON decode request body as %T", item),
					http.StatusBadRequest)
			}

			// Encoder for response.
			enc := json.NewEncoder(w)

			// Search the old stored item.
			// If not found, return 404.
			id := mux.Vars(r)["id"]
			stored := store.Get(item.GetType(), id).(storablePatchable)
			if stored == nil {
				w.WriteHeader(http.StatusNotFound)
				enc.Encode(map[string]interface{}{
					"status":  http.StatusNotFound,
					"message": "not found",
				})
				return
			}

			// Do the actual patching.
			stored.PatchWith(item)

			// Enforce ID, then put it back into the store.
			stored.SetID(id)
			store.Put(stored)

			w.WriteHeader(http.StatusOK)
			enc.Encode(map[string]interface{}{
				"status": http.StatusOK,
				nounp: []interface{}{
					stored,
				},
				noun: stored,
			})
		}

		Delete := func(w http.ResponseWriter, r *http.Request) {

			item := f.Make().(Storable)

			// encoder for returning protocol
			enc := json.NewEncoder(w)

			// search the old stored item
			// if not found, return 404
			id := mux.Vars(r)["id"]
			var ok bool
			if v := store.Get(item.GetType(), id); v == nil {
				w.WriteHeader(http.StatusNotFound)
				enc.Encode(map[string]interface{}{
					"status":  http.StatusNotFound,
					"message": "not found",
				})
				return
			} else if item, ok = v.(Storable); !ok {
				w.WriteHeader(http.StatusInternalServerError)
				enc.Encode(map[string]interface{}{
					"status":  http.StatusInternalServerError,
					"message": "internal server error",
				})
				return
			}

			store.Delete(item.GetType(), item.GetID())

			w.WriteHeader(http.StatusOK)
			enc.Encode(map[string]interface{}{
				"status": http.StatusOK,
				nounp: []interface{}{
					item,
				},
				noun: item,
			})

		}

		// path to use
		pathSingular := path.Join(base, noun, "{id}")
		pathPlural := path.Join(base, nounp)

		r := mux.NewRouter()
		r.HandleFunc(pathSingular, Retrieve).Methods("GET")
		r.HandleFunc(pathSingular, Update).Methods("PUT")
		r.HandleFunc(pathSingular, Patch).Methods("PATCH")
		r.HandleFunc(pathSingular, Delete).Methods("DELETE")
		r.HandleFunc(pathPlural, List).Methods("GET")
		r.HandleFunc(pathPlural, Create).Methods("POST")

		return r
	}

}
