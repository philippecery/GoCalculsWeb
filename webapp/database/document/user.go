package document

import (
	"crypto/hmac"
	hash "crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/philippecery/maths/webapp/config"
	"github.com/philippecery/maths/webapp/constant"
	"github.com/philippecery/maths/webapp/util"
)

// User document
type User struct {
	UserID         string
	EmailAddress   string
	Password       string
	FirstName      string
	LastName       string
	Role           constant.UserRole
	Status         constant.UserStatus
	LastConnection *time.Time
	Token          string
	Expires        *time.Time
}

// UnregisteredUser document
type UnregisteredUser struct {
	UserID  string
	Role    constant.UserRole
	Status  constant.UserStatus
	Token   string
	Expires *time.Time
}

// RegisteredUser document
type RegisteredUser struct {
	UserID         string
	EmailAddress   string
	Password       string
	FirstName      string
	LastName       string
	Role           constant.UserRole
	Status         constant.UserStatus
	LastConnection *time.Time
}

// Link returns the registration link
func (u *User) Link() string {
	return "/register?token=" + u.Token
}

// ActionToken generates and returns a unique ID to pass as a query parameter for CSRF protection.
func (u *User) ActionToken() string {
	salt := util.GenerateRandomBytes(32)
	mac := hmac.New(hash.New, []byte(config.Config.Keys.ActionToken))
	mac.Write([]byte(u.UserID))
	mac.Write(salt)
	token := make([]byte, 0)
	token = append(token, salt...)
	token = append(token, mac.Sum(nil)...)
	return base64.URLEncoding.EncodeToString(token)
}

// Enabled returns true if this user's status is Enabled.
func (u *User) Enabled() bool {
	return u.Status == constant.Enabled
}

// IsAdmin return true is this user's role is Admin
func (u *User) IsAdmin() bool {
	return u.Role == constant.Admin
}

// IsTeacher return true is this user's role is Teacher
func (u *User) IsTeacher() bool {
	return u.Role == constant.Teacher
}

// IsStudent return true is this user's role is Student
func (u *User) IsStudent() bool {
	return u.Role == constant.Student
}

// IsUnregistered return true is this user's status is Unregistered
func (u *User) IsUnregistered() bool {
	return u.Status == constant.Unregistered
}

// IsDisabled return true is this user's status is Disabled
func (u *User) IsDisabled() bool {
	return u.Status == constant.Disabled
}
