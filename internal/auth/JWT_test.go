package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	userId, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("unable to generate uuid for testing %v", err)
	}
	tokenSecret := "jumbalaya"
	expiresIn := 12 * time.Hour

	stringToTest, err := MakeJWT(userId, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("error creating JWT %v", err)
	}

	ID, err := ValidateJWT(stringToTest, tokenSecret)
	if err != nil {
		t.Errorf("unable to validate jwt %v", err)
	}

	if ID != userId {
		t.Errorf("id does not match %v", err)
	}
}
