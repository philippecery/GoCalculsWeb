package model

import (
	"encoding/base64"
	"time"

	"github.com/philippecery/maths/webapp/util"

	"github.com/philippecery/libs/bytes"
	"github.com/philippecery/libs/cipher"
	"github.com/philippecery/maths/webapp/config"
	"github.com/philippecery/maths/webapp/constant/team"
	"github.com/philippecery/maths/webapp/constant/user"
)

// PII is to be used for PII data
type PII struct {
	cipher string
}

type Team struct {
	TeamID   string
	Type     team.Type
	Name     string
	Language string
	Status   team.Status
}

// User model
type User struct {
	UserID          string
	Password        string
	PasswordDate    time.Time
	FullName        string
	Language        string
	EmailAddress    *PII
	EmailAddressTmp *PII
	Role            user.Role
	Status          user.Status
	LastConnection  time.Time
	Token           string
	Expires         time.Time
	FailedAttempts  int
	Teams           []*Team
}

// UserProfile model
type UserProfile struct {
	UserID         string
	EmailAddress   *PII
	FullName       string
	Language       string
	LastConnection time.Time
	Teams          []*Team
}

// RegisteredUser model
// Replaces existing User model when user registers
type RegisteredUser struct {
	UserID         string
	EmailAddress   *PII
	Password       string
	FullName       string
	Language       string
	Role           user.Role
	Status         user.Status
	LastConnection time.Time
	PasswordDate   time.Time
}

// Student model
// Updates an existing User model
type Student struct {
	UserID   string
	FullName string
	Language string
	Grade    *Grade
}

// Link returns the registration link
func (u *User) Link() string {
	return "https://" + config.Config.Hostname + "/register?token=" + u.Token
}

// Enabled returns true if this user's status is Enabled.
func (u *User) Enabled() bool {
	return u.Status == user.Enabled
}

// IsAdmin return true is this user's role is Admin
func (u *User) IsAdmin() bool {
	return u.Role == user.Admin
}

// IsTeacher return true is this user's role is Teacher
func (u *User) IsTeacher() bool {
	return u.Role == user.ParentOrTeacher
}

// IsParent return true is this user's role is Teacher
func (u *User) IsParent() bool {
	return u.Role == user.ParentOrTeacher
}

// IsStudent return true is this user's role is Student
func (u *User) IsStudent() bool {
	return u.Role == user.ChildOrStudent
}

// IsChild return true is this user's role is Student
func (u *User) IsChild() bool {
	return u.Role == user.ChildOrStudent
}

// IsUnregistered return true is this user's status is Unregistered
func (u *User) IsUnregistered() bool {
	return u.Status == user.Unregistered
}

// IsDisabled return true is this user's status is Disabled
func (u *User) IsDisabled() bool {
	return u.Status == user.Disabled
}

// ActionToken generates and returns a unique ID to pass as a query parameter for CSRF protection.
func (u *User) ActionToken() string {
	return actionToken(u.UserID)
}

// VerifyUserActionToken verifies the provided action token is valid for the provided user.
func VerifyUserActionToken(actionToken string, userID string) bool {
	return verifyActionToken(actionToken, userID)
}

// ActionToken generates and returns a unique ID to pass as a query parameter for CSRF protection.
func (s *Student) ActionToken() string {
	return actionToken(s.UserID, s.Grade.GradeID)
}

// VerifyStudentActionToken verifies the provided action token is valid for the provided student and grade.
func VerifyStudentActionToken(actionToken string, userID, gradeID string) bool {
	return verifyActionToken(actionToken, userID, gradeID)
}

// Protect encrypts the PII data
func Protect(pii string) (*PII, error) {
	var ppii string
	var err error
	if ppii, err = util.ProtectPII(pii); err == nil {
		return &PII{ppii}, nil
	}
	return nil, err
}

// Reveal decrypts and returns the PII data
func (p *PII) Reveal() string {
	var err error
	var piiKey, protectedPIIBytes, piiBytes []byte
	piiKey, err = base64.StdEncoding.DecodeString(config.Config.Keys.PII)
	if err == nil {
		protectedPIIBytes, err = bytes.Decode(string(p.cipher))
		if err == nil {
			piiBytes, _ = cipher.Decrypt(&piiKey, protectedPIIBytes)
			defer bytes.Clear(&piiBytes)
			return string(piiBytes)
		}
	}
	return ""
}
