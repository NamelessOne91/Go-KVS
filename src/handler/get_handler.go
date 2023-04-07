package handler

import (
	"errors"
	"net/http"

	"github.com/NamelessOne91/Go-KVS/store"
	"github.com/go-chi/chi/v5"
)

// KeyValueGetHandler is called with a GET request to /v1/{key}
func KeyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	key := string(chi.URLParam(r, "key"))

	value, err := store.Get(key)
	if errors.Is(err, store.ErrorNoSuchKey) {
		http.Error(w, err.Error(), http.StatusNotFound) // 404
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // 500
		return
	}

	w.Write([]byte(value))
}
