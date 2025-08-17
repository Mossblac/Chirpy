package ext

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/Mossblac/Chirpy/internal/database"
	"github.com/google/uuid"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	DB             *database.Queries
}

type Params struct {
	Text    string    `json:"body"`
	Email   string    `json:"email"`
	User_id uuid.UUID `json:"user_id"`
}

type ReturnValError struct {
	Error string `json:"error"`
}

type ReturnCleanBody struct {
	Cleaned_body string `json:"cleaned_body"`
}

type NewUser struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type NewChirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
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

func (cfg *ApiConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {
	access, err := DevAccess()
	if err != nil {
		w.Write([]byte(fmt.Sprintln("requires Dev permission")))
		w.WriteHeader(403)
	}
	if access {
		cfg.FileserverHits.Store(0)
		w.Write([]byte(fmt.Sprintln("hit count reset")))
		fmt.Printf("hit counter reset\n")

		err := cfg.DB.DeleteAllUsers(r.Context())
		if err != nil {
			somethingwentwrong := ReturnValError{Error: "Something went wrong"}
			WriteJSONResponse(w, 500, somethingwentwrong)
			template := "Error: %v"
			w.Write([]byte(fmt.Sprintf(template, err)))
		}
		w.Write([]byte(fmt.Sprintln("Users Reset")))
		fmt.Printf("users reset\n")
	}
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

func (cfg *ApiConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	emailToSet := Params{}
	err := decoder.Decode(&emailToSet)
	if err != nil {
		somethingwentwrong := ReturnValError{Error: "Something went wrong"}
		template := "Error: %v"
		w.Write([]byte(fmt.Sprintf(template, err)))
		WriteJSONResponse(w, 500, somethingwentwrong)
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), emailToSet.Email)
	if err != nil {
		somethingwentwrong := ReturnValError{Error: "Something went wrong"}
		WriteJSONResponse(w, 500, somethingwentwrong)
		template := "Error: %v"
		w.Write([]byte(fmt.Sprintf(template, err)))
		return
	}
	newuser := NewUser{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	WriteJSONResponse(w, 201, newuser)
}

func (cfg *ApiConfig) CreateChirpHandler(w http.ResponseWriter, r *http.Request) {
	validChirp, user_id := ValidateChirp(w, r)
	chirpParams := database.CreateChirpParams{
		Body:   validChirp,
		UserID: user_id,
	}

	chirp, err := cfg.DB.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		template := "Error: %v"
		w.Write([]byte(fmt.Sprintf(template, err)))
	}

	newchirp := NewChirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	WriteJSONResponse(w, 201, newchirp)
}

func (cfg *ApiConfig) GetChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.DB.GetChirps(r.Context())
	if err != nil {
		somethingwentwrong := ReturnValError{Error: "Something went wrong"}
		template := "Error: %v"
		w.Write([]byte(fmt.Sprintf(template, err)))
		WriteJSONResponse(w, 500, somethingwentwrong)
	}

	WriteJSONResponse(w, 200, chirps)
}
