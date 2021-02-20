package constant

// OperationStatus type
type OperationStatus uint8

// HomeworkSessionStatus type
type HomeworkSessionStatus uint8

// List of operation statuses
const (
	Todo OperationStatus = iota
	Wrong
	Good
)

// List of homework statuses
const (
	Cancel HomeworkSessionStatus = iota
	Failed
	Timeout
	Success
)

// Logo return the logo for this status
func (s HomeworkSessionStatus) Logo() string {
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
