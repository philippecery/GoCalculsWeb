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
