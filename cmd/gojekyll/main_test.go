package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuild(t *testing.T) {
	err := parseAndRun([]string{"build", "-s", "../../testdata/example", "-q"})
	require.NoError(t, err)
}
