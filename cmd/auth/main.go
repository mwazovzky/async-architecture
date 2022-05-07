package main

import (
	"context"
	"database/sql"
	"fmt"
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
	port = os.Getenv("PORT")
	allowedOrigin = os.Getenv("ALLOWED_ORIGIN")
}

func main() {
	db := connectDB()
	defer db.Close()

	router := mux.NewRouter()
	authHandlers := handlers.NewAuthHandlers(db)
	authApi := router.PathPrefix("/auth").Subrouter()
	authApi.HandleFunc("/register", authHandlers.Register).Methods(http.MethodPost)
	authApi.HandleFunc("/login", authHandlers.Login).Methods(http.MethodPost)
	authApi.HandleFunc("/logout", authHandlers.Logout).Methods(http.MethodPost)
	authApi.HandleFunc("/check", authHandlers.Check).Methods(http.MethodPost)
	authApi.HandleFunc("/user", authHandlers.UserIndex).Methods(http.MethodGet)
	authApi.HandleFunc("/user/{id:[0-9]+}", authHandlers.UserUpdate).Methods(http.MethodPatch)
	authApi.HandleFunc("/user/{id:[0-9]+}", authHandlers.UserDelete).Methods(http.MethodDelete)

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

func connectDB() *sql.DB {
	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")
	dbuser := os.Getenv("DB_USER")
	dbpassword := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbhost, dbport, dbuser, dbpassword, dbname)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to db")

	return db
}
