package auth

import (
	"fmt"

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

/*write some unit tests to verify
the functions are working*/
