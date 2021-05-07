package homework

// Status type
type Status uint8

// List of homework statuses
const (
	Draft Status = iota
	Online
	Archived
)

// Logo return the logo for this status
func (s Status) Logo() string {
	var logo string
	switch s {
	case Draft:
		logo = "remove"
	case Online:
		logo = "thumbs-down"
	case Archived:
		logo = "time"
	}
	return logo
}

// SessionStatus type
type SessionStatus uint8

// List of homework session statuses
const (
	Cancel SessionStatus = iota
	Failed
	Timeout
	Success
)

// Logo return the logo for this status
func (s SessionStatus) Logo() string {
	var logo string
	switch s {
	case Cancel:
		logo = "remove"
	case Failed:
		logo = "thumbs-down"
	case Timeout:
		logo = "time"
	case Success:
		logo = "thumbs-up"
	}
	return logo
}
