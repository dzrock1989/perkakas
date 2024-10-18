package perkakas

import "strings"

// IsEmpty validate value is empty
func IsEmpty(str string) bool {
	if len(str) <= 0 || str == "" {
		return true
	}

	return false
}

// IsEqual compare two string
func IsEqual(str1 string, str2 string) bool {
	return strings.Compare(str1, str2) == 0
}
