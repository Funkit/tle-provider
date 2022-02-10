package utils

// IsALetter Check if char is an ASCII letter or not
func IsALetter(element byte) bool {
	val := ((element >= 'a') && (element <= 'z')) || ((element >= 'A') && (element <= 'Z'))

	return val
}
