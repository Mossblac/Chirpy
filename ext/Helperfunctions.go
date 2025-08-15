package ext

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
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
