package util

import (
	"bytes"
	"crypto/rand"
	hash "crypto/sha256"
	"encoding/base64"
	"encoding/binary"

	memory "github.com/philippecery/libs/bytes"
	"golang.org/x/crypto/pbkdf2"
)

const (
	minIterations = 10000
	maxIterations = 15000 // Must be greater than minIterations and not greater than 65535
)

// ProtectPassword returns the Base64-encoded protected password
func ProtectPassword(password string) string {
	salt := make([]byte, hash.Size)
	rand.Read(salt)
	iter, _ := GetNumberInRange(minIterations, maxIterations)
	hashedPwd := hashPassword([]byte(password), salt, iter)
	defer memory.Clear(&hashedPwd, &salt)
	iterBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(iterBytes, uint16(iter))
	return base64.StdEncoding.EncodeToString(Concat(salt, iterBytes, hashedPwd))
}

// VerifyPassword verifies the submitted password against the actual one
func VerifyPassword(submitted, base64Blob string) bool {
	if blob, err := base64.StdEncoding.DecodeString(base64Blob); err == nil {
		iter := int(binary.BigEndian.Uint16(blob[hash.Size : hash.Size+2]))
		return bytes.Equal(hashPassword([]byte(submitted), blob[:hash.Size], iter), blob[hash.Size+2:])
	}
	return false
}

func hashPassword(password, salt []byte, iter int) []byte {
	return pbkdf2.Key(password, salt, iter, hash.Size, hash.New)
}
