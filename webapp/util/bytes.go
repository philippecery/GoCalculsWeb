package util

// Concat concatenates byte slices into a single slice
func Concat(bytesToConcat ...[]byte) []byte {
	bytes := make([]byte, 0)
	for _, b := range bytesToConcat {
		bytes = append(bytes, b...)
	}
	return bytes
}
