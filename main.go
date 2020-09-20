package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/nitinda/microservice-with-go/handlers"
)

func main() {
	l := log.New(os.Stdout, "product-api ", log.LstdFlags)

	// Create the handlers
	ph := handlers.NewProducts(l)

	// Create new server MUX and register the handlers
	sm := http.NewServeMux()
	sm.Handle("/", ph)

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
		l.Println("Starting Server on port 8000")
		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error Staring Server on port 8000\n%s", err)
			os.Exit(1)
		}
	}()

	// Create OS Signal Channel
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Recieved terminate, gracefulshutdown", sig)

	// Create new Context
	tc, _ := context.WithTimeout(context.Background(), 5*time.Second)

	// Graceful Shutdown
	s.Shutdown(tc)
}
