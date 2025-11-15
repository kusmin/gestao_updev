package testutil

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

// RandomPassword generates a random, URL-safe password-like string with at least 16 characters.
func RandomPassword() string {
	const size = 16
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		log.Printf("failed to generate random password: %v", err)
		return "fallback-" + base64.RawURLEncoding.EncodeToString([]byte("gestao-updev"))
	}
	return base64.RawURLEncoding.EncodeToString(b)
}
