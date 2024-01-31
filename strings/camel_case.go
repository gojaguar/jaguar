package strings

// CamelCase works identical to PascalCase with the difference that
// CamelCase will return the first value in the string as a lowercase character.
// In short, _my_field_name_2 becomes xMyFieldName_2.
func CamelCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'x')
		i++
	}
	// If the first letter is a lowercase, we keep it as is.
	if len(t) == 0 && isASCIILower(s[i]) {
		t = append(t, s[i])
		t, i = appendLowercaseSequence(s, i, t)
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		t, i = appendLowercaseSequence(s, i, t)
	}
	return string(t)
}
