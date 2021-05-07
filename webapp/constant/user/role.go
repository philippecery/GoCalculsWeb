package user

// Role type
type Role uint8

// List of roles
const (
	ChildOrStudent Role = iota
	ParentOrTeacher
	Principal
	Admin
)

var roles = []string{"Child/Student", "Parent/Teacher", "Principal", "Admin"}

func (r Role) String() string {
	if r.IsValid() {
		return roles[r]
	}
	return ""
}

// IsValid returns true if this role is valid
func (r Role) IsValid() bool {
	return int(r) < len(roles)
}
