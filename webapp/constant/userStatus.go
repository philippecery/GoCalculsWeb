package constant

// UserStatus type
type UserStatus uint8

// List of user status
const (
	Unregistered UserStatus = iota
	Disabled
	Enabled
)

var statuses = []string{"Unregistered", "Disabled", "Enabled"}

func (s UserStatus) String() string {
	if s.IsValid() {
		return statuses[s]
	}
	return ""
}

// IsValid returns true if this status is valid
func (s UserStatus) IsValid() bool {
	return s >= 0 && int(s) < len(statuses)
}
