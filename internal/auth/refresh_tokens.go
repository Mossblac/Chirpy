package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	rand.Read(key)

	HexString := hex.EncodeToString(key)

	return HexString, nil
}
