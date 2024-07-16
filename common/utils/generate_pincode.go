package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandomString(length int) (string, error) {
	// Create a byte slice with half the desired length
	bytes := make([]byte, length/2)

	// Read random bytes
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Encode the bytes to a hex string and truncate to the desired length
	randomString := hex.EncodeToString(bytes)[:length]
	return randomString, nil
}
