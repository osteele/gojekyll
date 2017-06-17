package gojekyll

// VariableMap is a map of strings to interface values, for use in template processing.
type VariableMap map[string]interface{}

// Bool returns m[k] if it's a bool; else defaultValue.
func (m VariableMap) Bool(k string, defaultValue bool) bool {
	if val, found := m[k]; found {
		if v, ok := val.(bool); ok {
			return v
		}
	}
	return defaultValue
}

// String returns m[k] if it's a string; else defaultValue.
func (m VariableMap) String(k string, defaultValue string) string {
	if val, found := m[k]; found {
		if v, ok := val.(string); ok {
			return v
		}
	}
	return defaultValue
}

// MergeVariableMaps creates a new variable map that merges its arguments,
// from first to last.
func MergeVariableMaps(ms ...VariableMap) VariableMap {
	result := VariableMap{}
	for _, m := range ms {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}
