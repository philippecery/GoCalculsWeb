package session

import (
	"strings"

	"github.com/philippecery/maths/webapp/constant"
)

// UserInformation contains only the user information we want to keep in the user session
type UserInformation struct {
	UserID    string
	FirstName string
	LastName  string
	Role      constant.UserRole
}

// IsAdmin returns true if this user's role is Admin
func (u *UserInformation) IsAdmin() bool {
	return u.HasRole(constant.Admin)
}

// IsTeacher returns true if this user's role is Teacher
func (u *UserInformation) IsTeacher() bool {
	return u.HasRole(constant.Teacher)
}

// IsStudent returns true if this user's role is Student
func (u *UserInformation) IsStudent() bool {
	return u.HasRole(constant.Student)
}

// HasRole returns true if this user has the specified role
func (u *UserInformation) HasRole(role constant.UserRole) bool {
	return u.Role == role
}

// CheckRoleByURI returns true if this user has the required role to access that URI
// A URI starting with:
//  - /admin requires Admin role
//  - /teacher requires Teacher role
//  - /student requires Student role
//  - /user is accessible to all authenticated users
func (u *UserInformation) CheckRoleByURI(uri string) bool {
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
