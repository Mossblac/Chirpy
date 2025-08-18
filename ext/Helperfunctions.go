package ext

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func WriteJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
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

func WordCleaner(bodytext string) string {
	words := strings.Split(bodytext, " ")
	for i := range words {
		if strings.ToLower(words[i]) ==
			"kerfuffle" || strings.ToLower(words[i]) ==
			"sharbert" || strings.ToLower(words[i]) ==
			"fornax" {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

func ValidateChirp(w http.ResponseWriter, r *http.Request) (string, uuid.UUID) {

	decoder := json.NewDecoder(r.Body)
	param := Params{}
	err := decoder.Decode(&param)
	if err != nil {
		WriteError(w, err)
		return "", uuid.Nil
	}
	if len(param.Text) <= 140 && len(param.Text) > 0 {
		cleantext := WordCleaner(param.Text)
		return cleantext, param.User_id

	} else if len(param.Text) > 140 {
		sendTooLong := ReturnValError{Error: "Chirp is too long"}
		WriteJSONResponse(w, 400, sendTooLong)

	} else {
		invalidChirp := ReturnValError{Error: "Chirp is invalid"}
		WriteJSONResponse(w, 400, invalidChirp)
	}
	return "", uuid.Nil
}

func WriteError(w http.ResponseWriter, err error) {
	somethingwentwrong := ReturnValError{Error: "Something went wrong\n\n"}
	template := "Error: %v\n\n"
	w.Write([]byte(fmt.Sprintf(template, err)))
	WriteJSONResponse(w, 500, somethingwentwrong)
}
