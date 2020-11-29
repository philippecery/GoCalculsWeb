package constant

// UserRole type
type UserRole uint8

// List of roles
const (
	None UserRole = iota
	Admin
	Teacher
	Student
)

var roles = []string{"None", "Admin", "Teacher", "Student"}

func (r UserRole) String() string {
	if r.IsValid() {
		return roles[r]
	}
	return ""
}

// IsValid returns true if this role is valid
func (r UserRole) IsValid() bool {
	return r >= 0 && int(r) < len(roles)
}
