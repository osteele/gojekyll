package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseAndRun(t *testing.T) {
	err := ParseAndRun([]string{"build", "-s", "../testdata/example", "-q"})
	require.NoError(t, err)
}
