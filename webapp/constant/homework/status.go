package homework

// SessionStatus type
type SessionStatus uint8

// List of homework statuses
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
