package handler

import (
	"io"
	"net/http"

	"github.com/NamelessOne91/Go-KVS/store"
	"github.com/go-chi/chi/v5"
)

// KeyValuePutHandler is called with a PUT request to /v1/{key}
// and expects the value of the key to be read from the HTTP Request's body
func KeyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	value, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // 500
		return
	}

	err = store.Put(key, string(value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // 500
		return
	}

	// HTTP 201
	w.WriteHeader(http.StatusCreated)
}
