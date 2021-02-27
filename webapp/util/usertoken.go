package util

import (
	"crypto/hmac"
	hash "crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/philippecery/maths/webapp/config"
)

func GenerateUserToken(userID string) (string, time.Time) {
	expirationTime := time.Now().Add(time.Hour * time.Duration(config.Config.UserTokenValidity))
	mac := hmac.New(hash.New, []byte(config.Config.Keys.UserToken))
	mac.Write([]byte(userID))
	mac.Write([]byte(fmt.Sprintf("%d", expirationTime.Unix())))
	return base64.URLEncoding.EncodeToString(mac.Sum(nil)), expirationTime
}
