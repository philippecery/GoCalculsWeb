package document

import (
	"bytes"
	"crypto/hmac"
	hash "crypto/sha256"
	"encoding/base64"

	"github.com/philippecery/libs/crng"

	"github.com/philippecery/maths/webapp/config"
)

// ActionToken generates and returns a unique ID to pass as a query parameter for CSRF protection.
func (u *User) ActionToken() string {
	salt, _ := crng.GetBytes(32)
	token := make([]byte, 0)
	token = append(token, salt...)
	token = append(token, generateActionToken(salt, u.UserID)...)
	return base64.URLEncoding.EncodeToString(token)
}

// VerifyUserActionToken verifies the provided action token is valid for the provided user.
func VerifyUserActionToken(actionToken string, userID string) bool {
	return verifyActionToken(actionToken, userID)
}

// ActionToken generates and returns a unique ID to pass as a query parameter for CSRF protection.
func (s *Student) ActionToken() string {
	salt, _ := crng.GetBytes(32)
	token := make([]byte, 0)
	token = append(token, salt...)
	token = append(token, generateActionToken(salt, s.UserID, s.GradeID)...)
	return base64.URLEncoding.EncodeToString(token)
}

// VerifyStudentActionToken verifies the provided action token is valid for the provided student and grade.
func VerifyStudentActionToken(actionToken string, userID, gradeID string) bool {
	return verifyActionToken(actionToken, userID, gradeID)
}

// ActionToken generates and returns a unique ID to pass as a query parameter for CSRF protection.
func (g *Grade) ActionToken() string {
	salt, _ := crng.GetBytes(32)
	token := make([]byte, 0)
	token = append(token, salt...)
	token = append(token, generateActionToken(salt, g.GradeID)...)
	return base64.URLEncoding.EncodeToString(token)
}

// VerifyGradeActionToken verifies the provided action token is valid for the provided grade.
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
