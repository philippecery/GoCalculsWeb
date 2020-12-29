package document

import (
	"crypto/hmac"
	hash "crypto/sha256"
	"encoding/base64"

	"github.com/philippecery/maths/webapp/config"
	"github.com/philippecery/maths/webapp/util"
)

// Grade document
type Grade struct {
	GradeID     string
	Name        string
	Description string
	MentalMath  *Homework
	ColumnForm  *Homework
}

// Homework sub-document
type Homework struct {
	NbAdditions       int
	NbSubstractions   int
	NbMultiplications int
	NbDivisions       int
	Time              int
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
