package tags

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var argTests = []struct {
	in          string
	optionCount int
	positional  []string
}{
	{`filename`, 0, []string{"filename"}},
	{`filename a=1`, 1, []string{"filename"}},
	{`filename a=1 b=2`, 2, []string{"filename"}},
	{`filename a='1' b=2`, 2, []string{"filename"}},
	{`filename a='1 b=' c`, 1, []string{"filename", "c"}},
	{`a=1 b=2`, 2, []string{}},
	{`a='1' b=2`, 2, []string{}},
	{`arg1 arg2`, 0, []string{"arg1", "arg2"}},
}

func TestFilters(t *testing.T) {
	for i, test := range argTests {
		t.Run(fmt.Sprintf("%02d", i+1), func(t *testing.T) {
			actual, err := ParseArgs(test.in)
			require.NoError(t, err)
			require.Equal(t, test.optionCount, len(actual.Options), "options count in %q -> #%v", test.in, actual)
			require.Equal(t, test.positional, actual.Args, "args in %q -> #%v", test.in, actual)
		})
	}
}
