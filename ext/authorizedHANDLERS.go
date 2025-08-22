package ext

import (
	"database/sql"
	"encoding/hex"
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
	Text     string    `json:"body"`
	Email    string    `json:"email"`
	User_id  uuid.UUID `json:"user_id"`
	Password string    `json:"password"`
}

type ReturnValError struct {
	Error string `json:"error"`
}

type ReturnCleanBody struct {
	Cleaned_body string `json:"cleaned_body"`
}

type NewUser struct {
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Email         string    `json:"email"`
	Token         string    `json:"token"`
	RefreshToken  string    `json:"refresh_token"`
	Is_Chirpy_Red bool      `json:"is_chirpy_red"`
}

type NewChirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type HookEvent struct {
	Event string `json:"event"`
	Data  UData  `json:"data"`
}

type UData struct {
	User_id string `json:"user_id"`
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

func (cfg *ApiConfig) CreateChirpHandler(w http.ResponseWriter, r *http.Request) {
	user_id, err := VerifyFromAccessTokenHeader(cfg, w, r)
	if err != nil {
		w.WriteHeader(401)
		return
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
		return
	}

	newchirp := NewChirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	WriteJSONResponse(w, 201, newchirp)
	fmt.Printf("chirp created by user: %v\n", newchirp.ID)
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

	token, err := auth.MakeJWT(user.ID, Sword, time.Duration(1)*time.Hour)
	if err != nil {
		WriteError(w, err)
		return
	}

	refresh, err := auth.MakeRefreshToken()
	if err != nil {
		WriteError(w, err)
		return
	}

	RefreshParam := database.CreateRefreshTokenParams{
		Token:  refresh,
		UserID: user.ID,
	}

	cfg.DB.CreateRefreshToken(r.Context(), RefreshParam)

	displayUser := NewUser{
		ID:            user.ID,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		Email:         user.Email,
		Token:         token,
		RefreshToken:  refresh,
		Is_Chirpy_Red: user.IsChirpyRed,
	}

	WriteJSONResponse(w, 200, displayUser)
	fmt.Println("User Logged in")

}

func (cfg *ApiConfig) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(401)
		WriteError(w, err)
		return
	}

	Rfesh, err := cfg.DB.GetRefresh(r.Context(), refreshToken)
	if err != nil {
		WriteError(w, err)
		return
	}
	currentTime := time.Now()

	if currentTime.After(Rfesh.ExpiresAt) || Rfesh.RevokedAt.Valid {
		w.WriteHeader(401)
		return
	} else {
		Sword := cfg.SecretWord
		Updatedtoken, err := auth.MakeJWT(Rfesh.UserID, Sword, time.Duration(1)*time.Hour)
		if err != nil {
			WriteError(w, err)
			return
		}

		type NewToken struct {
			Token string `json:"token"`
		}
		newToken := NewToken{Token: Updatedtoken}
		WriteJSONResponse(w, 200, newToken)
		fmt.Printf("Token refreshed: %v\n", newToken.Token)
	}
}

func (cfg *ApiConfig) RevokeHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(401)
		WriteError(w, err)
		return
	}

	Rfesh, err := cfg.DB.GetRefresh(r.Context(), refreshToken)
	if err != nil {
		WriteError(w, err)
		return
	}

	cfg.DB.RevokeRefresh(r.Context(), Rfesh.Token)
	w.WriteHeader(204)
	fmt.Printf("token revoked\n")
}

func (cfg *ApiConfig) ResetEmailAndPasswordHandler(w http.ResponseWriter, r *http.Request) {
	user_id, err := VerifyFromAccessTokenHeader(cfg, w, r)
	if err != nil {
		WriteError(w, err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	loginparam := Params{}
	err = decoder.Decode(&loginparam)
	if err != nil {
		WriteError(w, err)
		return
	}
	HexPass := hex.EncodeToString([]byte(loginparam.Password))
	resetParams := database.ResetEmailAndPasswordParams{
		Email:          loginparam.Email,
		HashedPassword: HexPass,
		ID:             user_id,
	}

	resetuser, err := cfg.DB.ResetEmailAndPassword(r.Context(), resetParams)
	if err != nil {
		WriteError(w, err)
		return
	}
	currentTime := time.Now()
	displayRUser := NewUser{
		ID:            resetuser.ID,
		CreatedAt:     resetuser.CreatedAt,
		UpdatedAt:     currentTime,
		Email:         resetuser.Email,
		Is_Chirpy_Red: resetuser.IsChirpyRed,
	}
	WriteJSONResponse(w, 200, displayRUser)
	fmt.Println("email and password reset success")
}

func (cfg *ApiConfig) DeleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	user_id, err := VerifyFromAccessTokenHeader(cfg, w, r)
	if err != nil {
		WriteError(w, err)
		return
	}
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

	if user_id == chirp.UserID {
		err := cfg.DB.DeleteChirp(r.Context(), chirp.ID)
		if err != nil {
			WriteError(w, err)
			return
		} else {
			w.WriteHeader(204)
			fmt.Println("deleted success")
			return
		}
	} else {
		w.WriteHeader(403)
		fmt.Println("chirp userId does not match user ID")
		return
	}
}

func (cfg *ApiConfig) UpgradeToChirpyRedHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	WebHookEvent := HookEvent{}
	err := decoder.Decode(&WebHookEvent)
	if err != nil {
		WriteError(w, err)
		return
	}

	if WebHookEvent.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	} else {
		userID, err := uuid.Parse(WebHookEvent.Data.User_id)
		if err != nil {
			WriteError(w, err)
			return
		}
		err = cfg.DB.UpgradeToChirpyRed(r.Context(), userID)
		if err != nil {
			w.WriteHeader(404)
			fmt.Printf("%v", err)
			return
		} else {
			w.WriteHeader(204)
			fmt.Println("user upgrade success")
			return
		}
	}
}
