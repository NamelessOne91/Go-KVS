package main

import (
	"log"
	"net/http"

	"github.com/NamelessOne91/Go-KVS/handler"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {
		r.Get("/{key}", handler.KeyValueGetHandler)
		r.Put("/{key}", handler.KeyValuePutHandler)
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
