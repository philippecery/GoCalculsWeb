package util

import (
	crand "crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	mrand "math/rand"
)

// GenerateRandomBytes generates a random slice of bytes
func GenerateRandomBytes(length int) []byte {
	b := make([]byte, 32)
	if _, err := io.ReadFull(crand.Reader, b); err != nil {
		panic("Could not generate random bytes")
	}
	return b
}

// GenerateRandomBytesToBase64 returns a base64-encoded random slice of bytes
func GenerateRandomBytesToBase64(length int) string {
	return base64.URLEncoding.EncodeToString(GenerateRandomBytes(length))
}

var rng *mrand.Rand

func init() {
	var seed int64
	binary.Read(crand.Reader, binary.BigEndian, &seed)
	rng = mrand.New(mrand.NewSource(seed))
}

// GetNumber selects a random integer between 0 and the given maximum number.
func GetNumber(max int) (int, error) {
	if max > 0 {
		return rng.Intn(max), nil
	}
	return 0, fmt.Errorf("Invalid parameter 'max'")
}

// GetNumberInRange selects a random integer within the given range.
// Returns an error if range is invalid
func GetNumberInRange(min, max int) (int, error) {
	if min >= 0 && min < max {
		result, _ := GetNumber(max)
		return result%(max-min) + min, nil
	}
	return 0, fmt.Errorf("Invalid parameters")
}
