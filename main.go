package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Server represents our HTTP server with its configurations and dependencies.
type Server struct {
	Addr    string
	Handler http.Handler
	Logger  *log.Logger
	// Add other dependencies or configurations here, e.g., database connections
	// DB *sql.DB
}

// NewServer creates and returns a new Server instance.
func NewServer(addr string, handler http.Handler, logger *log.Logger) *Server {
	return &Server{
		Addr:    addr,
		Handler: handler,
		Logger:  logger,
	}
}

// Start initiates the HTTP server.
func (s *Server) Start() error {
	s.Logger.Printf("Starting server on %s\n", s.Addr)
	srv := &http.Server{
		Addr:    s.Addr,
		Handler: s.Handler,
		// Optional: Add timeouts for robustness
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	return srv.ListenAndServe()
}

func main() {
	// Create a simple HTTP handler
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "404 Not found")
	})

	// Create a logger
	logger := log.Default()

	// Create a new server instance
	myServer := NewServer(":8080", mux, logger)

	// Start the server
	if err := myServer.Start(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Server failed to start: %v", err)
	}
}
