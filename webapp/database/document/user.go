package document

import (
	"time"

	"github.com/philippecery/maths/webapp/config"
	"github.com/philippecery/maths/webapp/constant"
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
	LastConnection time.Time
	Token          string
	Expires        time.Time
	GradeID        string
	FailedAttempts int
	PasswordDate   time.Time
}

// UserProfile document
type UserProfile struct {
	UserID         string
	EmailAddress   string
	FirstName      string
	LastName       string
	LastConnection time.Time
}

// UnregisteredUser document
// Used to create a new User document when admin creates new user
type UnregisteredUser struct {
	UserID  string
	Role    constant.UserRole
	Status  constant.UserStatus
	Token   string
	Expires time.Time
}

// RegisteredUser document
// Replaces existing User document when user registers
type RegisteredUser struct {
	UserID         string
	EmailAddress   string
	Password       string
	FirstName      string
	LastName       string
	Role           constant.UserRole
	Status         constant.UserStatus
	LastConnection time.Time
	PasswordDate   time.Time
}

// Student document
// Updates an existing User document
type Student struct {
	UserID    string
	FirstName string
	LastName  string
	GradeID   string
	Grade     *Grade
}

// Link returns the registration link
func (u *User) Link() string {
	return "/register?token=" + u.Token
}

// Link returns the registration link
func (u *UnregisteredUser) Link() string {
	return "https://" + config.Config.Hostname + "/register?token=" + u.Token
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
