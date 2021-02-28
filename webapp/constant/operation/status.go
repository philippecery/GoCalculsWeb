package operation

// Status type
type Status uint8

// List of operation statuses
const (
	Todo Status = iota
	Wrong
	Good
)
