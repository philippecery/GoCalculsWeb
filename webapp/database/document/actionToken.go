package document

import (
	"bytes"
	"crypto/hmac"
	hash "crypto/sha256"
	"encoding/base64"

	"github.com/philippecery/maths/webapp/config"
	"github.com/philippecery/maths/webapp/util"
)

// ActionToken generates and returns a unique ID to pass as a query parameter for CSRF protection.
func (u *User) ActionToken() string {
	salt := util.GenerateRandomBytes(32)
	token := make([]byte, 0)
	token = append(token, salt...)
	token = append(token, generateActionToken(salt, u.UserID)...)
	return base64.URLEncoding.EncodeToString(token)
}

func VerifyUserActionToken(actionToken string, userID string) bool {
	return verifyActionToken(actionToken, userID)
}

// ActionToken generates and returns a unique ID to pass as a query parameter for CSRF protection.
func (s *Student) ActionToken() string {
	salt := util.GenerateRandomBytes(32)
	token := make([]byte, 0)
	token = append(token, salt...)
	token = append(token, generateActionToken(salt, s.UserID, s.GradeID)...)
	return base64.URLEncoding.EncodeToString(token)
}

func VerifyStudentActionToken(actionToken string, userID, gradeID string) bool {
	return verifyActionToken(actionToken, userID, gradeID)
}

// ActionToken generates and returns a unique ID to pass as a query parameter for CSRF protection.
func (g *Grade) ActionToken() string {
	salt := util.GenerateRandomBytes(32)
	mac := hmac.New(hash.New, []byte(config.Config.Keys.ActionToken))
	mac.Write([]byte(g.GradeID))
	mac.Write(salt)
	token := make([]byte, 0)
	token = append(token, salt...)
	token = append(token, mac.Sum(nil)...)
	return base64.URLEncoding.EncodeToString(token)
}

func VerifyGradeActionToken(actionToken string, gradeID string) bool {
	return verifyActionToken(actionToken, gradeID)
}

func verifyActionToken(actionToken string, ids ...string) bool {
	if token, err := base64.URLEncoding.DecodeString(actionToken); err == nil {
		return bytes.Equal(token[32:], generateActionToken(token[:32], ids...))
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
