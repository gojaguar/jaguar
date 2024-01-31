package strings

// Is c an ASCII lower-case letter?
func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// appendLowercaseSequence appends the lowercase sequence from s that begins at i into t
// returns the new t that contains all the chain of characters that should be lowercase
// and the new index where to start counting from.
func appendLowercaseSequence(s string, i int, t []byte) ([]byte, int) {
	for i+1 < len(s) && isASCIILower(s[i+1]) {
		i++
		t = append(t, s[i])
	}
	return t, i
}
