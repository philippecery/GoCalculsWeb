package session

import "github.com/philippecery/maths/webapp/constant"

// UserInformation contains only the user information we want to keep in the user session
type UserInformation struct {
	UserID    string
	FirstName string
	LastName  string
	Role      constant.UserRole
}

// IsAdmin returns true is this user's role is Admin
func (u *UserInformation) IsAdmin() bool {
	return u.Role == constant.Admin
}

// IsTeacher returns true is this user's role is Teacher
func (u *UserInformation) IsTeacher() bool {
	return u.Role == constant.Teacher
}

// IsStudent returns true is this user's role is Student
func (u *UserInformation) IsStudent() bool {
	return u.Role == constant.Student
}
