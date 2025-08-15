package main

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/Mossblac/Chirpy/ext"
)

func main() {
	cfg := ext.ApiConfig{
		FileserverHits: atomic.Int32{},
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/healthz", ext.HealthzHandler)

	mux.HandleFunc("GET /admin/metrics", cfg.ShowCountHandler)

	mux.HandleFunc("POST /admin/reset", cfg.ResetCountHandler)

	mux.HandleFunc("POST /api/validate_chirp", ext.ValidateChirpHandler)

	mux.Handle("/app/", cfg.MetricsINC(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))

	/*
		MetricsINC returns a handlerfunc that matches the Handler interface. it runs the code within that function just like
		ResetCounterHandler and ShowCountHandler, but then runs the handler that is set as its input

		so here it provides a function that adds to the Hitcounter, and then after it has done that returns the handler
		that is set as the input, which in this case is the fileserver handler

		so instead of a handler function assigned to its own page, this function "wraps" around
		the function for the main page: "/app/"
	*/

	err := http.ListenAndServe(":8080", mux)
	fmt.Printf("%v", err)

}
