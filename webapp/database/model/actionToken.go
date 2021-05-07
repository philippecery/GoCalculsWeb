package model

import (
	b "bytes"
	"crypto/hmac"
	hash "crypto/sha256"
	"encoding/base64"

	"github.com/philippecery/libs/bytes"
	"github.com/philippecery/libs/crng"

	"github.com/philippecery/maths/webapp/config"
)

func actionToken(ids ...string) string {
	salt, _ := crng.GetBytes(32)
	token := make([]byte, 0)
	token = append(token, salt...)
	token = append(token, generateActionToken(salt, ids...)...)
	return bytes.Encode(token)
}

func verifyActionToken(actionToken string, ids ...string) bool {
	if token, err := base64.URLEncoding.DecodeString(actionToken); err == nil {
		return b.Equal(token[32:], generateActionToken(token[:32], ids...))
	}
	return false
}

func generateActionToken(salt []byte, ids ...string) []byte {
	mac := hmac.New(hash.New, []byte(config.Config.Keys.ActionToken))
	for _, id := range ids {
		mac.Write([]byte(id))
	}
	mac.Write(salt)
	return mac.Sum(nil)
}
