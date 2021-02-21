package util

import (
	"encoding/base64"

	"github.com/philippecery/libs/bytes"

	"github.com/philippecery/libs/cipher"
	"github.com/philippecery/libs/hmac"
	"github.com/philippecery/maths/webapp/config"
)

// ProtectUserID protects the email address used as a user identifier
func ProtectUserID(userID string) (string, error) {
	macKey, err := base64.StdEncoding.DecodeString(config.Config.Keys.UserID)
	if err == nil {
		return base64.StdEncoding.EncodeToString(hmac.Generate(&macKey, []byte(userID))), nil
	}
	return "", err
}

// ProtectPII protects PII data at rest
func ProtectPII(pii string) (string, error) {
	var err error
	var piiKey, ciphertext []byte
	piiBytes := []byte(pii)
	piiKey, err = base64.StdEncoding.DecodeString(config.Config.Keys.PII)
	if err == nil {
		ciphertext, err = cipher.Encrypt(&piiKey, &piiBytes)
		return base64.StdEncoding.EncodeToString(ciphertext), nil
	}
	return "", err
}

func UnprotectPII(protectedPII string) (string, error) {
	var err error
	var piiKey, protectedPIIBytes, piiBytes []byte
	piiKey, err = base64.StdEncoding.DecodeString(config.Config.Keys.PII)
	if err == nil {
		protectedPIIBytes, err = base64.StdEncoding.DecodeString(protectedPII)
		if err == nil {
			piiBytes, err = cipher.Decrypt(&piiKey, protectedPIIBytes)
			defer bytes.Clear(&piiBytes)
			return string(piiBytes), nil
		}
	}
	return "", err
}
