package helpers

// MergeStringMaps creates a new variable map that merges its arguments,
// from first to last.
func MergeStringMaps(ms ...map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	for _, m := range ms {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}
