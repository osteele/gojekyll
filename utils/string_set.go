package utils

// A StringSet is a set of strings, represented as a map.
type StringSet map[string]bool

// MakeStringSet creates a characteristic function map that tests for presence in an array.
func MakeStringSet(a []string) StringSet {
	set := map[string]bool{}
	for _, s := range a {
		set[s] = true
	}
	return set
}

// AddStrings modifies the set to include the strings in an array.
func (ss StringSet) AddStrings(a []string) {
	for _, s := range a {
		ss[s] = true
	}
}
