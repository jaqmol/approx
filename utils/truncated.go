package utils

// Truncated ...
func Truncated(input []byte, byteLength int) (output []byte) {
	if len(input) <= byteLength {
		return input
	}
	headLen := (byteLength / 2) - 2
	tailLen := byteLength - headLen - 1
	output = make([]byte, 0)
	output = append(output, input[:headLen]...)
	output = append(output, []byte("...")...)
	output = append(output, input[len(input)-tailLen:]...)
	return
}
