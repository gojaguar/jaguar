package strings

// CamelCase is a special case of PascalCase with the difference that
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
	return string(append(t, lookupAndReplacePascalCaseWords(s, i)...))
}
