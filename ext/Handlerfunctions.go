package ext

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
}

func (cfg *ApiConfig) MetricsINC(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})

}

func (cfg *ApiConfig) ShowCountHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("Hits: %v", cfg.FileserverHits.Load())))
}

func (cfg *ApiConfig) ResetCountHandler(w http.ResponseWriter, r *http.Request) {
	cfg.FileserverHits.Store(0)
	w.Write([]byte(fmt.Sprintln("hit count reset")))
}

func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // write the content-type header

	w.WriteHeader(200) // write the status code

	message := []byte("OK")
	n, err := w.Write(message) // write the body text
	if err != nil {
		fmt.Printf("Error writing resposne %v\n", err)
		return
	}
	fmt.Printf("wrote %d bytes to response\n", n)
}
