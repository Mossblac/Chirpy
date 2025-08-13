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

	mux.HandleFunc("/healthz", ext.HealthzHandler)

	mux.HandleFunc("/metrics", cfg.ShowCountHandler)

	mux.HandleFunc("/reset", cfg.ResetCountHandler)

	mux.Handle("/app/", cfg.MetricsINC(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))

	err := http.ListenAndServe(":8080", mux)
	fmt.Printf("%v", err)

}
