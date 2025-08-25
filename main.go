package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/Mossblac/Chirpy/ext"
	_ "github.com/lib/pq"
)

func main() {
	dbQueries, err := ext.DatabaseAccess()
	if err != nil {
		log.Fatal(err)
	}

	Sword, Pkey, err := ext.GetENVWords()
	if err != nil {
		log.Fatal(err)
	}

	cfg := ext.ApiConfig{
		FileserverHits: atomic.Int32{},
		DB:             dbQueries,
		SecretWord:     Sword,
		PolkaKey:       Pkey,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/healthz", ext.HealthzHandler)

	mux.HandleFunc("GET /admin/metrics", cfg.ShowCountHandler)

	mux.HandleFunc("POST /admin/reset", cfg.ResetHandler)

	mux.HandleFunc("POST /api/chirps", cfg.CreateChirpHandler)

	mux.HandleFunc("POST /api/users", cfg.CreateUserHandler)

	mux.HandleFunc("GET /api/chirps", cfg.GetChirpsHandler)

	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.GetSingleChirpHandler)

	mux.HandleFunc("POST /api/login", cfg.UserLoginHandler)

	mux.HandleFunc("POST /api/refresh", cfg.RefreshHandler)

	mux.HandleFunc("POST /api/revoke", cfg.RevokeHandler)

	mux.HandleFunc("PUT /api/users", cfg.ResetEmailAndPasswordHandler)

	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.DeleteChirpHandler)

	mux.HandleFunc("POST /api/polka/webhooks", cfg.UpgradeToChirpyRedHandler)

	mux.Handle("/app/", cfg.MetricsINC(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))

	err = http.ListenAndServe(":8080", mux)
	fmt.Printf("%v", err)

}
