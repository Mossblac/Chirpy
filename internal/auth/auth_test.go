package auth

import (
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
