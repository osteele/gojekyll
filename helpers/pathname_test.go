package helpers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTrimExt(t *testing.T) {
	require.Equal(t, "/a/b", TrimExt("/a/b.c"))
	require.Equal(t, "/a/b", TrimExt("/a/b"))
}
