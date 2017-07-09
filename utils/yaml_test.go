package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const mapYaml = "a: 1\nb: 2"
const listYaml = "- a\n- b"

func TestUnmarshalYAML(t *testing.T) {
	var d interface{}
	err := UnmarshalYAMLInterface([]byte(mapYaml), &d)
	require.NoError(t, err)
	switch d := d.(type) {
	case map[interface{}]interface{}:
		require.Len(t, d, 2)
		require.Equal(t, 1, d["a"])
	default:
		require.IsType(t, d, map[interface{}]interface{}{})
	}

	err = UnmarshalYAMLInterface([]byte(listYaml), &d)
	require.NoError(t, err)
	require.IsType(t, d, []interface{}{})
	switch d := d.(type) {
	case []interface{}:
		require.Len(t, d, 2)
		require.Equal(t, "a", d[0])
	default:
		require.IsType(t, d, map[interface{}]interface{}{})
	}
}
