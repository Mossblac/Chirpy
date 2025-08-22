package ext

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mossblac/Chirpy/internal/auth"
	"github.com/Mossblac/Chirpy/internal/database"
)

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
