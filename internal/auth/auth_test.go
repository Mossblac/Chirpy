package auth

import (
	"net/http"
	"testing"
)

func TestAuthorization(t *testing.T) {
	password := "$&jibbajabba6987"
	hash, _ := HashPassword(password)
	err := CheckPasswordHash(password, hash)
	result := err
	if result != nil {
		t.Errorf("result returned: %v expected nil", err)
	}

}

func TestMismatch(t *testing.T) {
	password := "&*io224543hjhjh"
	hash, _ := HashPassword(password)
	password2 := "thisshouldfail"
	err := CheckPasswordHash(password2, hash)
	result := err
	if result == nil {
		t.Errorf("password mismatch, should error with: %v", err)
	}

}

func TestHashing(t *testing.T) {
	password := "nonsense4765!!!!"
	_, err := HashPassword(password)
	result := err
	if result != nil {
		t.Errorf("hashing function returned error: %v", err)
	}

}

func TestGetBearerRegular(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer my_test_token_123")

	token, err := GetBearerToken(headers)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if token != "my_test_token_123" {
		t.Errorf("Expected 'my_test_token_123', got: '%s'", token)
	}
}

func TestGetBearerEmpty(t *testing.T) {
	headers := http.Header{}
	headers.Set("notAuthorized", "This is not an authorization")

	_, err := GetBearerToken(headers)

	if err == nil {
		t.Errorf("invalid format, should have errored")
	}
}

func TestGetBearerMissingToken(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer")

	_, err := GetBearerToken(headers)

	if err == nil {
		t.Errorf("no token present, should have errored")
	}
}
