package handler

import (
	"net/http"

	"github.com/NamelessOne91/Go-KVS/store"
	"github.com/NamelessOne91/Go-KVS/transaction"
	"github.com/go-chi/chi/v5"
)

// KeyValueDeleteHandler is called with a DELETE request to /v1/{key}
func KeyValueDeleteHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	err := store.Delete(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // 500
		return
	}

	transaction.Logger.WriteDelete(key)
	w.WriteHeader(http.StatusOK)
}
