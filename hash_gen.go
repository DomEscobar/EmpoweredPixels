package main

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
)

const (
	passwordIterations = 10000
	passwordKeyLength  = 32
)

func main() {
	password := "AgentPower123!"
	salt := "AGENT_TESTER_SALT_BASE64_=="

	saltBytes, _ := base64.StdEncoding.DecodeString(salt)
	hash := pbkdf2.Key([]byte(password), saltBytes, passwordIterations, passwordKeyLength, sha512.New)

	fmt.Printf("Salt: %s\n", salt)
	fmt.Printf("Hash: %s\n", base64.StdEncoding.EncodeToString(hash))
}
