package ext

import "testing"

func TestWordCleaner(t *testing.T) {
	input := "I hear Mastodon is better than Chirpy. sharbert I need to migrate"
	result := WordCleaner(input)
	expected := "I hear Mastodon is better than Chirpy. **** I need to migrate"
	if result != expected {
		t.Errorf("WordCleaner(input) returned %v, expected %v", result, expected)
	}
}

func TestWordCleanerALLBADWORDS(t *testing.T) {
	input := "!!! kerfuffle sharbert fornax !$!"
	result := WordCleaner(input)
	expected := "!!! **** **** **** !$!"
	if result != expected {
		t.Errorf("WordCleaner(input) returned %v, expected %v", result, expected)
	}
}

func TestWordCleanerAlreadyClean(t *testing.T) {
	input := "No bad words heRe"
	result := WordCleaner(input)
	expected := "No bad words heRe"
	if result != expected {
		t.Errorf("WordCleaner(input) returned %v, expected %v", result, expected)
	}
}
