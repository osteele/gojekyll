package utils

// StringList adds methods to []string
type StringList []string

// Reject returns a copy of StringList without elements that its argument tests true on.
func (sl StringList) Reject(p func(string) bool) StringList {
	result := make([]string, 0, len(sl))
	for _, s := range sl {
		if !p(s) {
			result = append(result, s)
		}
	}
	return result
}

// SearchStrings returns a bool indicating whether array contains the string.
// Unlike sort.SearchStrings, it does not require a sorted array.
// This is useful when the array length is low; else consider sorting it and using
// sort.SearchStrings, or creating a map[string]bool.
func SearchStrings(array []string, s string) bool {
	for _, item := range array {
		if item == s {
			return true
		}
	}
	return false
}
