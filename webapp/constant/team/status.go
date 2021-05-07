package team

// Status type
type Status uint8

// List of org status
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
	return int(s) < len(statuses)
}
