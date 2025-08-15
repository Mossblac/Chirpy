package ext

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
}

func (cfg *ApiConfig) MetricsINC(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		fmt.Printf("hit added\n")
		next.ServeHTTP(w, r)
	})

}

/* the middleware function above, returns a handler function like those below, it follows the
signature of the Handler interface and therefore returns a Handler.
the "next" input handler is the handler that the middleware wraps around.
see the use case in main()*/

func (cfg *ApiConfig) ShowCountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	template := "<html> <body> <h1>Welcome, Chirpy Admin</h1> <p>Chirpy has been visited %d times!</p> </body> </html>"

	w.Write([]byte(fmt.Sprintf(template, cfg.FileserverHits.Load())))
	fmt.Printf("Total Hits: %v\n", cfg.FileserverHits.Load())
}

func (cfg *ApiConfig) ResetCountHandler(w http.ResponseWriter, r *http.Request) {
	cfg.FileserverHits.Store(0)
	w.Write([]byte(fmt.Sprintln("hit count reset")))
	fmt.Printf("hit counter reset\n")
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

func ValidateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Text string `json:"body"`
	}

	type returnValError struct {
		Error string `json:"error"`
	}

	type returnValValid struct {
		Valid bool `json:"valid"`
	}

	somethingwentwrong := returnValError{
		Error: "Something went wrong",
	}

	decoder := json.NewDecoder(r.Body)
	param := params{}
	err := decoder.Decode(&param)
	if err != nil {
		writeJSONResponse(w, 500, somethingwentwrong)
		return
	}
	if len(param.Text) <= 140 && len(param.Text) > 0 {
		sendValid := returnValValid{
			Valid: true,
		}
		writeJSONResponse(w, 200, sendValid)
	} else if len(param.Text) > 140 {
		sendTooLong := returnValError{
			Error: "Chirp is too long",
		}
		writeJSONResponse(w, 400, sendTooLong)
	} else {
		invalidChirp := returnValError{
			Error: "Chirp is invalid",
		}
		writeJSONResponse(w, 400, invalidChirp)
	}
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	dat, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(dat)
}
