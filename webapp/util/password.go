package util

import (
	"bytes"
	"crypto/rand"
	hash "crypto/sha256"
	"encoding/base64"
)

// ProtectPassword returns the BASE64-encoded protected password
func ProtectPassword(password string) string {
	salt := make([]byte, 32)
	rand.Read(salt)
	h := hash.New()
	h.Write(salt)
	h.Write([]byte(password))
	hashedPwd := make([]byte, 0)
	hashedPwd = append(hashedPwd, salt...)
	hashedPwd = append(hashedPwd, h.Sum(nil)...)
	return base64.StdEncoding.EncodeToString(hashedPwd)
}

// VerifyPassword verifies the submitted password against the actual one
func VerifyPassword(submitted, expected string) bool {
	if hashedPwd, err := base64.StdEncoding.DecodeString(expected); err == nil {
		h := hash.New()
		h.Write(hashedPwd[:32])
		h.Write([]byte(submitted))
		return bytes.Equal(h.Sum(nil), hashedPwd[32:])
	}
	return false
}
