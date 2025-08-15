package ext

import (
	"encoding/json"
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

	type Params struct {
		Text string `json:"body"`
	}

	type ReturnValError struct {
		Error string `json:"error"`
	}

	type ReturnCleanBody struct {
		Cleaned_body string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	param := Params{}
	err := decoder.Decode(&param)
	if err != nil {
		somethingwentwrong := ReturnValError{Error: "Something went wrong"}
		WriteJSONResponse(w, 500, somethingwentwrong)
		return
	}
	if len(param.Text) <= 140 && len(param.Text) > 0 {
		cleantext := WordCleaner(param.Text)
		sendCleanText := ReturnCleanBody{Cleaned_body: cleantext}
		WriteJSONResponse(w, 200, sendCleanText)

	} else if len(param.Text) > 140 {
		sendTooLong := ReturnValError{Error: "Chirp is too long"}
		WriteJSONResponse(w, 400, sendTooLong)

	} else {
		invalidChirp := ReturnValError{Error: "Chirp is invalid"}
		WriteJSONResponse(w, 400, invalidChirp)
	}
}
