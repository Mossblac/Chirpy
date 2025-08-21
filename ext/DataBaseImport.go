package ext

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Mossblac/Chirpy/internal/database"
	"github.com/joho/godotenv"
)

func DatabaseAccess() (*database.Queries, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DB_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	dbQueries := database.New(db)

	return dbQueries, nil
}

func DevAccess() (access bool, err error) {
	err = godotenv.Load()
	if err != nil {
		return false, fmt.Errorf("error loading .env file: %v", err)
	}

	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		return false, nil
	}

	return true, nil
}

func GetSecretWord() (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", fmt.Errorf("error loading .env file: %v", err)
	}

	word := os.Getenv("SECRETWORD")

	return word, nil
}
