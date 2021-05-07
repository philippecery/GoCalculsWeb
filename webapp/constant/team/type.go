package team

// Type type
type Type uint8

// List of org types
const (
	Family Type = iota + 1
	School
)

var types = []string{"Family", "School"}

func (t Type) String() string {
	if t.IsValid() {
		return types[t]
	}
	return ""
}

// IsValid returns true if this type is valid
func (t Type) IsValid() bool {
	return int(t) < len(types)
}
