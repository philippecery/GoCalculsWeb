package team

// Role type
type Role uint8

// List of roles
const (
	Normal Role = iota
	Admin
	Root
)

var roles = []string{"Normal", "Admin", "Root"}

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
