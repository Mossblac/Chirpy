package auth

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	passwordinput := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(passwordinput, 10)
	if err != nil {
		return "", fmt.Errorf("error creating hash: %v", err)
	}
	hashstring := string(hash)
	return hashstring, nil
}

func CheckPasswordHash(password, hash string) error {
	hashinput := []byte(hash)
	passwordinput := []byte(password)
	err := bcrypt.CompareHashAndPassword(hashinput, passwordinput)
	if err != nil {
		return fmt.Errorf("error when comparing hash to password: %v", err)
	}
	return nil
}

func GetBearerToken(headers http.Header) (string, error) {
	rawBearer := headers.Get("Authorization")
	if rawBearer == "" {
		return "", fmt.Errorf("no Authorization header found")
	}

	SplitBearer := strings.Split(rawBearer, " ")

	if SplitBearer[0] != "Bearer" {
		return "", fmt.Errorf("bearer token not formatted correctly")
	} else if SplitBearer[0] == "Bearer" && len(SplitBearer) == 1 {
		return "", fmt.Errorf("bearer token missing")
	} else {
		rawBearer := strings.TrimSpace(strings.TrimPrefix(headers.Get("Authorization"), "Bearer"))
		return rawBearer, nil
	}

}

func GetAPIKey(headers http.Header) (string, error) {
	rawAPIkey := headers.Get("Authorization")
	if rawAPIkey == "" {
		return "", fmt.Errorf("no Authorization header found")
	}

	SplitAPIKey := strings.Split(rawAPIkey, " ")

	if SplitAPIKey[0] != "ApiKey" {
		return "", fmt.Errorf("API key not formatted correctly")
	} else if SplitAPIKey[0] == "ApiKey" && len(SplitAPIKey) == 1 {
		return "", fmt.Errorf("API key missing")
	} else {
		APIkey := strings.TrimSpace(strings.TrimPrefix(headers.Get("Authorization"), "ApiKey"))
		return APIkey, nil
	}
}
