package util

// Clear zeroizes the provided slice of bytes, then set the slice to nil
func Clear(bytes *[]byte) {
	if bytes != nil {
		for i := range *bytes {
			(*bytes)[i] = 0
		}
		*bytes = nil
	}
}
