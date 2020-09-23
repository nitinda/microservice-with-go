package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"

	"github.com/nitinda/microservice-with-go/data"
	"github.com/nitinda/microservice-with-go/handlers"
)

func main() {
	l := log.New(os.Stdout, "products-api ", log.LstdFlags)
	v := data.NewValidation()

	// create the handlers
	ph := handlers.NewProducts(l, v)

	// create a new serve mux and register the handlers
	sm := mux.NewRouter()

	// handlers for API
	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/api/products", ph.ListAll)
	getR.HandleFunc("/api/products/{id:[0-9]+}", ph.ListSingle)

	putR := sm.Methods(http.MethodPut).Subrouter()
	putR.HandleFunc("/api/products", ph.Update)
	putR.Use(ph.MiddlewareValidateProduct)

	postR := sm.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/api/products", ph.Create)
	postR.Use(ph.MiddlewareValidateProduct)

	deleteR := sm.Methods(http.MethodDelete).Subrouter()
	deleteR.HandleFunc("api/products/{id:[0-9]+}", ph.Delete)

	// handler for documentation
	opts := middleware.RedocOpts{
		SpecURL: "/swagger.yaml",
		Path:    "/api/docs",
	}
	sh := middleware.Redoc(opts, nil)

	getR.Handle("/api/docs", sh)
	getR.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// Create new Server
	s := &http.Server{
		Addr:         ":8000",           // Condifgure the bind address
		Handler:      sm,                // Set the Default Handler
		ErrorLog:     l,                 // ErrorLog specifies an optional logger for errors accepting
		ReadTimeout:  5 * time.Second,   // ReadTimeout is the maximum duration for reading the entire request
		WriteTimeout: 10 * time.Second,  // WriteTimeout is the maximum duration before timing out
		IdleTimeout:  120 * time.Second, // IdleTimeout is the maximum amount of time to wait for the
	}

	// Start the Server
	go func() {
		l.Println("Starting server on port 8000")

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// Block until a signal is received.
	sig := <-sigChan
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	tc, _ := context.WithTimeout(context.Background(), 5*time.Second)
	s.Shutdown(tc)
}
