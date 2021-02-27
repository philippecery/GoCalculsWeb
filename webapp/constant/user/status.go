package user

// Status type
type Status uint8

// MaxFailedAttempts is the number of times the user is allowed to enter a wrong password before getting his/her account disabled
const MaxFailedAttempts = 5

// List of user status
const (
	Unregistered Status = iota
	Disabled
	Enabled
)

var statuses = []string{"Unregistered", "Disabled", "Enabled"}

func (s Status) String() string {
	if s.IsValid() {
		return statuses[s]
	}
	return ""
}

// IsValid returns true if this status is valid
func (s Status) IsValid() bool {
	return s >= 0 && int(s) < len(statuses)
}
