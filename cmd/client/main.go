package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"async-architecure/http/handlers"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var port string
var allowedOrigin string

func init() {
	godotenv.Load()
	port = os.Getenv("CLIENT_PORT")
	allowedOrigin = os.Getenv("CLIENT_ALLOWED_ORIGIN")
}

func main() {
	router := mux.NewRouter()
	clientHandlers := handlers.NewClientHandlers()
	clientApi := router.PathPrefix("/client").Subrouter()
	clientApi.HandleFunc("/ping", clientHandlers.Ping).Methods(http.MethodGet)

	// CORS
	cors := gohandlers.CORS(
		gohandlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		gohandlers.AllowedOrigins([]string{allowedOrigin}),
		gohandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
	)

	server := &http.Server{
		Addr:         port,
		Handler:      cors(router),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		log.Println("Starting http server at", port)
		err := server.ListenAndServe()
		if err != nil {
			log.Println("Error", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	sig := <-sigChan
	log.Printf("Recieved terminate signal, graceful shutdown, signal: [%s]", sig)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	server.Shutdown(ctx)
}
