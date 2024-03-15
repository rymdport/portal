package convert

// ToNullTerminated connverts a regular string into a null terminated byte string.
func ToNullTerminated(input string) []byte {
	terminated := make([]byte, len(input)+1)
	copy(terminated, input)
	return terminated
}
