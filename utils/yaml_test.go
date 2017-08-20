package utils

import (
	"testing"

	yaml "gopkg.in/yaml.v2"

	"github.com/stretchr/testify/require"
)

const mapYaml = "a: 1\nb: 2"
const listYaml = "- a\n- b"

func TestUnmarshalYAML(t *testing.T) {
	var d interface{}
	err := UnmarshalYAMLInterface([]byte(mapYaml), &d)
	require.NoError(t, err)
	switch d := d.(type) {
	case yaml.MapSlice:
		require.Len(t, d, 2)
		require.Equal(t, yaml.MapItem{Key: "a", Value: 1}, d[0])
	default:
		require.IsType(t, map[interface{}]interface{}{}, d)
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
