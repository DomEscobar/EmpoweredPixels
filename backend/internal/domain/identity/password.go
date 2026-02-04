package identity

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"

	"golang.org/x/crypto/pbkdf2"
)

const (
	passwordIterations = 10000
	passwordKeyLength  = 32
)

func GenerateSalt() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}

func HashPassword(password string, base64Salt string) (string, error) {
	saltBytes, err := base64.StdEncoding.DecodeString(base64Salt)
	if err != nil {
		return "", err
	}

	hash := pbkdf2.Key([]byte(password), saltBytes, passwordIterations, passwordKeyLength, sha512.New)
	return base64.StdEncoding.EncodeToString(hash), nil
}

func PasswordMatches(password string, base64Salt string, expectedHash string) bool {
	hash, err := HashPassword(password, base64Salt)
	if err != nil {
		return false
	}

	return hash == expectedHash
}
