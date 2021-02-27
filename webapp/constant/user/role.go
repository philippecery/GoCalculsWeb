package user

// Role type
type Role uint8

// List of roles
const (
	Child Role = iota
	Student
	Parent
	Teacher
	Principal
	Admin
)

var roles = []string{"User", "Student", "Parent", "Teacher", "Principal", "Admin"}

func (r Role) String() string {
	if r.IsValid() {
		return roles[r]
	}
	return ""
}

// IsValid returns true if this role is valid
func (r Role) IsValid() bool {
	return r >= 0 && int(r) < len(roles)
}
