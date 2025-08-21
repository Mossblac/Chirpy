package ext

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/Mossblac/Chirpy/internal/auth"
	"github.com/Mossblac/Chirpy/internal/database"

	"github.com/google/uuid"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	DB             *database.Queries
	SecretWord     string
}

type Params struct {
	Text               string    `json:"body"`
	Email              string    `json:"email"`
	User_id            uuid.UUID `json:"user_id"`
	Password           string    `json:"password"`
	Expires_in_seconds int       `json:"expiration"`
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
	Token     string    `json:"token"`
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
			WriteError(w, err)
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
	paramUnhashed := Params{}
	err := decoder.Decode(&paramUnhashed)
	if err != nil {
		WriteError(w, err)
		return
	}

	if paramUnhashed.Password == "unset" || paramUnhashed.Password == "" {
		WritePasswordError(w, err)
	} else {

		hash, err := auth.HashPassword(paramUnhashed.Password)
		if err != nil {
			WriteError(w, err)
			return
		}

		hashedparams := database.CreateUserParams{
			Email:          paramUnhashed.Email,
			HashedPassword: hash,
		}

		user, err := cfg.DB.CreateUser(r.Context(), hashedparams)
		if err != nil {
			WriteError(w, err)
			return
		}
		newuser := NewUser{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		}

		WriteJSONResponse(w, 201, newuser)
		fmt.Println("User created")
	}
}

func (cfg *ApiConfig) CreateChirpHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		WriteError(w, err)
		w.WriteHeader(401)
	}

	user_id, err := auth.ValidateJWT(token, cfg.SecretWord)
	if err != nil {
		WriteError(w, err)
		w.WriteHeader(401)
	}

	validChirp, _ := ValidateChirp(w, r)
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
		WriteError(w, err)
	}

	apiChirps := make([]NewChirp, 0, len(chirps))

	for _, c := range chirps {
		apiChirps = append(apiChirps, NewChirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		})
	}

	WriteJSONResponse(w, 200, apiChirps)
}

func (cfg *ApiConfig) GetSingleChirpHandler(w http.ResponseWriter, r *http.Request) {
	inputid := (r).PathValue("chirpID")
	chirpid, err := uuid.Parse(inputid)
	if err != nil {
		WriteError(w, err)
		return
	}
	chirp, err := cfg.DB.GetSingleChirpByID(r.Context(), chirpid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(404)
			return
		} else {
			WriteError(w, err)
			return
		}
	}
	newchirp := NewChirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	WriteJSONResponse(w, 200, newchirp)
}

func (cfg *ApiConfig) UserLoginHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	loginparam := Params{}
	err := decoder.Decode(&loginparam)
	if err != nil {
		WriteError(w, err)
		return
	}

	user, err := cfg.DB.GetUserViaEmail(r.Context(), loginparam.Email)
	if err != nil {
		w.WriteHeader(401)
		notAuthorized := "incorrect email or password"
		w.Write([]byte(notAuthorized))
		return
	}

	err = auth.CheckPasswordHash(loginparam.Password, user.HashedPassword)
	if err != nil {
		w.WriteHeader(401)
		notAuthorized := "incorrect email or password"
		w.Write([]byte(notAuthorized))
		return
	}

	Sword := cfg.SecretWord

	var duration time.Duration

	if loginparam.Expires_in_seconds > 0 && loginparam.Expires_in_seconds < 3600 {
		duration = time.Duration(loginparam.Expires_in_seconds) * time.Second
	} else {
		duration = time.Duration(1) * time.Hour
	}

	token, err := auth.MakeJWT(user.ID, Sword, duration)
	if err != nil {
		WriteError(w, err)
	}

	displayUser := NewUser{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	}

	WriteJSONResponse(w, 200, displayUser)
	fmt.Println("User Logged in")

}
