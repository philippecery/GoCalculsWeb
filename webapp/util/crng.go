package util

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

// GenerateRandomBytes generates a random slice of bytes
func GenerateRandomBytes(length int) []byte {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic("Could not generate random bytes")
	}
	return b
}

// GenerateRandomBytesToBase64 returns a base64-encoded random slice of bytes
func GenerateRandomBytesToBase64(length int) string {
	return base64.URLEncoding.EncodeToString(GenerateRandomBytes(length))
}
