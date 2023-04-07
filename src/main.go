package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NamelessOne91/Go-KVS/handler"
	"github.com/NamelessOne91/Go-KVS/middleware"
	"github.com/NamelessOne91/Go-KVS/transaction"
	"github.com/go-chi/chi/v5"
	chimid "github.com/go-chi/chi/v5/middleware"
)

func service() http.Handler {
	r := chi.NewRouter()
	r.Use(chimid.Recoverer)

	fmt.Println("Server running at 127.0.0.1:8080")
	r.Route("/v1", func(r chi.Router) {
		r.Get("/{key}", handler.KeyValueGetHandler)
		r.With(middleware.Limiter).Put("/{key}", handler.KeyValuePutHandler)
		r.Delete("/{key}", handler.KeyValueDeleteHandler)
	})

	return r
}

func main() {
	transaction.InitTransactionLogger()

	// The HTTP Server
	server := &http.Server{Addr: "127.0.0.1:8080", Handler: service()}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig
		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
				cancel()
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
