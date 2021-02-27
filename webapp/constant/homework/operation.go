package homework

// OperationStatus type
type OperationStatus uint8

// List of operation statuses
const (
	Todo OperationStatus = iota
	Wrong
	Good
)
