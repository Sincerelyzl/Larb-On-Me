package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

const (
	saltLength = 16
	timeHash   = 1
	memory     = 64 * 1024
	threads    = 4
	keyLength  = 32
)

func generateSalt() ([]byte, error) {
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func HashPassword(password string) (string, error) {
	salt, err := generateSalt()
	if err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, timeHash, memory, threads, keyLength)
	return base64.RawStdEncoding.EncodeToString(salt) + "$" + base64.RawStdEncoding.EncodeToString(hash), nil
}

func VerifyPassword(password, hashedPassword string) bool {
	parts := split(hashedPassword, '$')
	if len(parts) != 2 {
		return false
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}
	hash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	newHash := argon2.IDKey([]byte(password), salt, timeHash, memory, threads, keyLength)
	return subtle.ConstantTimeCompare(hash, newHash) == 1
}

func split(s string, sep byte) []string {
	n := 1
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			n++
		}
	}
	a := make([]string, n)
	n--
	i := 0
	for i < n {
		m := i + 1
		for m < len(s) && s[m] != sep {
			m++
		}
		a[i] = s[i:m]
		i = m + 1
	}
	a[n] = s[i:]
	return a
}
