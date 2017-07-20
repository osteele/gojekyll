package cache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithFile(t *testing.T) {
	require.NoError(t, resetCache())
	callCount := 0

	stringMaker := func() (string, error) {
		callCount++
		return "ok", nil
	}
	errMaker := func() (string, error) {
		return "", fmt.Errorf("expected error")
	}
	requireOk := func(t *testing.T, s string, err error, expectCallCount int) {
		require.Equal(t, expectCallCount, callCount)
		require.NoError(t, err)
		require.Equal(t, "ok", s)
	}

	t.Run("calls underlying function", func(t *testing.T) {
		s, err := WithFile("h1", "c1", stringMaker)
		requireOk(t, s, err, 1)
	})

	t.Run("cache hit when keys match", func(t *testing.T) {
		callCount = 0
		s, err := WithFile("h1", "c1", stringMaker)
		requireOk(t, s, err, 0)
	})

	t.Run("cache miss when header differs", func(t *testing.T) {
		callCount = 0
		s, err := WithFile("h2", "c1", stringMaker)
		requireOk(t, s, err, 1)
	})

	t.Run("cache miss when content differs", func(t *testing.T) {
		callCount = 0
		s, err := WithFile("h1", "c2", stringMaker)
		requireOk(t, s, err, 1)
	})

	t.Run("propagates error", func(t *testing.T) {
		_, err := WithFile("h1-err", "c1", errMaker)
		require.Error(t, err)
		require.Contains(t, err.Error(), "expected error")
	})
}
