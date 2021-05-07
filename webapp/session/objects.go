package session

import (
	"strings"
	"time"

	"github.com/philippecery/maths/webapp/constant/user"
)

// UserInformation contains only the user information we want to keep in the user session
type UserInformation struct {
	UserID         string
	Name           string
	Language       string
	Role           user.Role
	LastConnection time.Time
	TeamID          string
}

// IsAdmin returns true if this user's role is Admin
func (u *UserInformation) IsAdmin() bool {
	return u.HasRole(user.Admin)
}

// IsTeacher returns true if this user's role is Teacher
func (u *UserInformation) IsTeacher() bool {
	return u.HasRole(user.ParentOrTeacher)
}

// IsStudent returns true if this user's role is Student
func (u *UserInformation) IsStudent() bool {
	return u.HasRole(user.ChildOrStudent)
}

// IsParent returns true if this user's role is Admin
func (u *UserInformation) IsParent() bool {
	return u.HasRole(user.ParentOrTeacher)
}

// IsChild returns true if this user's role is Admin
func (u *UserInformation) IsChild() bool {
	return u.HasRole(user.ChildOrStudent)
}

// HasRole returns true if this user has the specified role
func (u *UserInformation) HasRole(role user.Role) bool {
	return u.Role == role
}

// HasAccessTo returns true if this user has the required role to access that URI
// A URI starting with:
//  - /admin requires Admin role
//  - /teacher requires Teacher role
//  - /student requires Student role
//  - /user is accessible to all authenticated users
func (u *UserInformation) HasAccessTo(uri string) bool {
	path := strings.Split(uri, "/")
	if len(path) > 1 {
		switch path[1] {
		case "admin":
			return u.IsAdmin()
		case "teacher":
			return u.IsTeacher()
		case "student":
			return u.IsStudent()
		case "user":
			return true
		}
	}
	return false
}
